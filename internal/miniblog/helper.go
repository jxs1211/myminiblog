// Copyright 2022 Innkeeper Jayflow <jxs121@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package miniblog

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marmotedu/miniblog/internal/miniblog/store"
	"github.com/marmotedu/miniblog/internal/pkg/log"
	"github.com/marmotedu/miniblog/pkg/db"
)

const (
	// recommendedHomeDir defines the default directory where miniblog service configuration is placed.
	recommendedHomeDir = ".miniblog"

	// defaultConfigName specifies the default configuration file name for miniblog service.
	defaultConfigName = "miniblog.yaml"
)

// initConfig sets the configuration file name to be read, environment variables, and reads the configuration file content into viper.
func initConfig() {
	if cfgFile != "" {
		// Read from the configuration file specified in the command-line options.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find the user's home directory.
		home, err := os.UserHomeDir()
		// If failed to get the user's home directory, print `'Error: xxx` and exit the program (exit code 1).
		cobra.CheckErr(err)

		// Add the directory `$HOME/<recommendedHomeDir>` to the configuration file search paths.
		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))

		// Add the current directory to the configuration file search paths.
		viper.AddConfigPath(".")

		// Set the configuration file format to YAML (YAML format is clear and readable, and supports complex configuration structures).
		viper.SetConfigType("yaml")

		// Configuration file name (without file extension).
		viper.SetConfigName(defaultConfigName)
	}

	// Read matching environment variables.
	viper.AutomaticEnv()

	// Read environment variables with the prefix MINIBLOG, automatically converting lowercase to uppercase if it is miniblog.
	viper.SetEnvPrefix("MINIBLOG")

	// The following 2 lines replace '.' and '-' in the key string of viper.Get(key) with '_'.
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Read the configuration file. If the configuration file name is specified, use the specified configuration file; otherwise, search in the registered search paths.
	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}

	// Print the configuration file currently used by viper for debugging purposes.
	log.Debugw("Using config file", "file", viper.ConfigFileUsed())
}

// logOptions reads the log configuration from viper, builds `*log.Options`, and returns it.
// Note: The key in `viper.Get<Type>()` needs to be separated by `.` to match the same indentation as YAML.
func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

// initStore reads the db configuration, creates a gorm.DB instance, and initializes the miniblog store layer.
func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	ins, err := db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	_ = store.NewStore(ins)

	return nil
}
