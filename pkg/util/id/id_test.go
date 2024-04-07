// Copyright 2022 Innkeeper Jayflow <jxs121@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/marmotedu/miniblog.

package id

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenShortID(t *testing.T) {
	shortID := GenShortID()
	assert.NotEqual(t, "", shortID)
	assert.Equal(t, 6, len(shortID))
}

func BenchmarkGenShortID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenShortID()
	}
}

func BenchmarkGenShortIDTimeConsuming(b *testing.B) {
	b.StopTimer() // 调用该函数停止压力测试的时间计数

	shortId := GenShortID()
	if shortId == "" {
		b.Error("Failed to generate short id")
	}

	b.StartTimer() // 重新开始时间

	for i := 0; i < b.N; i++ {
		GenShortID()
	}
}

func workIn2Minute() {
	log.Println("starting to work")
	time.Sleep(10 * time.Second)
}

func doWork(ctx context.Context, d time.Duration) {
	// now := time.Now()
	delay := time.NewTicker(d)

	for {
		workIn2Minute()
		delay.Reset(d)

		// delay.Stop()
		select {
		case <-ctx.Done():
			delay.Stop()
		case <-delay.C:
		}
	}
}

func TestOverlappingTickerTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go doWork(ctx, time.Second*2)
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM) // Register for SIGTERM signal

	<-sigChan
	cancel()
}
