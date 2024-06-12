// Copyright 2022 Innkeeper Jayflow <jxs121@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// _ "go.uber.org/automaxprocs"

	"github.com/marmotedu/miniblog/internal/miniblog"
	"github.com/marmotedu/miniblog/pkg/util/tips"
)

func main() {
	// start()
	context, cancel := context.WithCancel(context.Background())
	go tips.Job3(context, 5*time.Second)
	// time.Sleep(60 * time.Second)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	<-sig
	cancel()
	log.Println("exit")
}

func start() {
	command := miniblog.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
