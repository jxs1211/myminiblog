package tips

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func workIn2Minute() {
	fmt.Println("do work in 2 minute")
	time.Sleep(2 * time.Minute)
}

func doWork(ctx context.Context, d time.Duration) {
	// now := time.Now()
	delay := time.NewTicker(d)

	for {
		// delay.Stop()
		workIn2Minute()
		// delay.Reset(d)

		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			delay.Stop()
		case <-delay.C:
		}
	}
}

func doWork2(ctx context.Context, d time.Duration) {
	// now := time.Now()
	delay := time.NewTicker(d)

	for {
		select {
		case <-ctx.Done():
			delay.Stop()
		case <-delay.C:
		}

		delay.Stop()
		workIn2Minute()
		delay.Reset(d)
	}
}

func TestOverlappingTickerTask(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go doWork(ctx, time.Second*2)
	// ch := make(chan struct{})
	sigal := make(chan os.Signal, 1)
	// add singal
	signal.Notify(sigal, syscall.SIGTERM)
	<-sigal
	cancel()
	t.Log("done task")
}

// Handle Errors of Deferred Calls to Prevent Silent Failures
func doSomething() error {
	file, err := os.Open("file.txt")
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	// read the file content and print it out
	fmt.Println(file.Name())
	fmt.Println(file.Stat())
	fmt.Println(file.Fd())

	fmt.Println()
	return err
}

// Always Keep Track of Your Goroutine's Lifetime
func Job(d time.Duration) {
	for ; ; time.Sleep(d) {
		fmt.Println("do something in job")
	}
}

func Job2(d time.Duration) {
	for {
		time.Sleep(d)
	}
}

func Job3(ctx context.Context, d time.Duration) {
	for {
		select {
		case <-ctx.Done():
			log.Println("cancel job")
			return
		default:
			// ...
			log.Println("start to work...")
			// time.Sleep(d)
			Sleep(ctx, d)
			log.Println("stop working...")
		}
	}
}

func Sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		log.Println("sleep: cancel job")
		return
	case <-t.C:
	}
}
