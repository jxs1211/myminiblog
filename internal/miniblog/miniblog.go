// Copyright 2022 Innkeeper Jayflow <jxs121@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package miniblog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/marmotedu/miniblog/internal/miniblog/controller/v1/user"
	"github.com/marmotedu/miniblog/internal/miniblog/store"
	"github.com/marmotedu/miniblog/internal/pkg/known"
	"github.com/marmotedu/miniblog/internal/pkg/log"
	mw "github.com/marmotedu/miniblog/internal/pkg/middleware"
	pb "github.com/marmotedu/miniblog/pkg/proto/miniblog/v1"
	"github.com/marmotedu/miniblog/pkg/token"
	"github.com/marmotedu/miniblog/pkg/version/verflag"
)

var cfgFile string

// NewMiniBlogCommand creates a *cobra.Command object. Afterwards, the Execute method of the Command object can be used to start the application.
func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		// Specify the name of the command, which will appear in the help information.
		Use: "miniblog",
		// A brief description of the command.
		Short: "A good Go practical project",
		// A detailed description of the command.
		Long: `A good Go practical project, used to create user with basic information.

Find more miniblog information at:
	https://github.com/marmotedu/miniblog#readme`,

		// When there is an error with the command, do not print the help information. Setting it to true allows for easy visibility of the error message when a command encounters an error.
		SilenceUsage: true,
		// Specify the Run function that will be executed when calling cmd.Execute(). If the function execution fails, it will return an error message.
		RunE: func(cmd *cobra.Command, args []string) error {
			// If `--version=true`, print the version and exit.
			verflag.PrintAndExitIfRequested()

			// Initialize logging.
			log.Init(logOptions())
			defer log.Sync() // flushes the log entries from the cache to the disk file.

			return run()
		},
		// configure the command to run without requiring command-line arguments.
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	// Run the function of configuration initialization when each command's
	// Execute method is called.
	cobra.OnInitialize(initConfig)

	// Define flags and configuration settings here.

	// Cobra supports persistent flags, which can be used by the assigned command and each sub-command under that command.
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the miniblog configuration file. Empty string for no configuration file.")

	// Cobra also supports local flags, which can only be used on the command it is bound to.
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add the --version flag.
	verflag.AddFlags(cmd.PersistentFlags())

	return cmd
}

// The run function is the actual entry point for the business logic.
func run() error {
	// Initialize the store layer
	if err := initStore(); err != nil {
		return err
	}

	// Set the signing key for the token package, used for token signing and parsing
	token.Init(viper.GetString("jwt-secret"), known.XUsernameKey)

	// Set Gin mode
	gin.SetMode(viper.GetString("runmode"))

	// Create a Gin engine
	g := gin.New()

	// Middleware functions for Gin: gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()}

	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return err
	}

	// Create and run an HTTP server
	httpsrv := startInsecureServer(g)

	// Create and run an HTTPS server
	httpssrv := startSecureServer(g)

	// Create and run a gRPC server
	grpcsrv := startGRPCServer()

	// Wait for an interrupt signal to gracefully shut down the server (with a 10-second timeout).
	quit := make(chan os.Signal, 1)
	// The kill command sends the syscall.SIGTERM signal by default
	// kill -2 sends the syscall.SIGINT signal, which is triggered by pressing CTRL + C
	// kill -9 sends the syscall.SIGKILL signal, but it can't be caught, so we don't need to handle it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // This won't block
	<-quit                                               // Block here, and only continue when one of the above signals is received
	log.Infow("Shutting down server ...")

	// Create a context to notify the server goroutine, giving it 10 seconds to complete the current requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shut down the server within 10 seconds (by completing the ongoing requests before shutting down)
	// If it takes more than 10 seconds, the server will time out and exit
	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}
	if err := httpssrv.Shutdown(ctx); err != nil {
		log.Errorw("Secure Server forced to shutdown", "err", err)
		return err
	}

	grpcsrv.GracefulStop()

	log.Infow("Server exiting")

	return nil
}

// startInsecureServer creates and runs an HTTP server.
func startInsecureServer(g *gin.Engine) *http.Server {
	// Create an instance of HTTP Server
	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	// Run the HTTP server. Start the server in a goroutine, so it doesn't block the normal shutdown process below.
	// Print a log message to indicate that the HTTP service is up and running, for troubleshooting purposes.
	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	return httpsrv
}

// startSecureServer creates and runs an HTTPS server.
func startSecureServer(g *gin.Engine) *http.Server {
	// Create an instance of HTTPS Server
	httpssrv := &http.Server{Addr: viper.GetString("tls.addr"), Handler: g}

	// Run the HTTPS server. Start the server in a goroutine, so it doesn't block the normal shutdown process below.
	// Print a log message to indicate that the HTTPS service is up and running, for troubleshooting purposes.
	log.Infow("Start to listening the incoming requests on https address", "addr", viper.GetString("tls.addr"))
	cert, key := viper.GetString("tls.cert"), viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			if err := httpssrv.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalw(err.Error())
			}
		}()
	}

	return httpssrv
}

// startGRPCServer creates and runs a gRPC server.
func startGRPCServer() *grpc.Server {
	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
	}

	// Create an instance of GRPC Server
	grpcsrv := grpc.NewServer()
	pb.RegisterMiniBlogServer(grpcsrv, user.New(store.S, nil))

	// Run the GRPC server. Start the server in a goroutine, so it doesn't block the normal shutdown process below.
	// Print a log message to indicate that the GRPC service is up and running, for troubleshooting purposes.
	log.Infow("Start to listening the incoming requests on grpc address", "addr", viper.GetString("grpc.addr"))
	go func() {
		if err := grpcsrv.Serve(lis); err != nil {
			log.Fatalw(err.Error())
		}
	}()

	return grpcsrv
}
