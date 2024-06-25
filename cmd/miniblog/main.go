// Copyright 2022 Innkeeper Jayflow <jxs121@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package main

import (
	"fmt"
	"os"
	"runtime"

	// _ "go.uber.org/automaxprocs"

	"github.com/marmotedu/miniblog/internal/miniblog"
)

func main() {
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	fmt.Println(runtime.NumCPU())
	start()
}

func start() {
	command := miniblog.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
