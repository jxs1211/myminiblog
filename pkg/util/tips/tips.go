package tips

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
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

//nolint:unused
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

func FetchExpensiveData() (int64, error) {
	fmt.Println("FetchExpensiveData called", time.Now())
	time.Sleep(3 * time.Second)
	return time.Now().Unix(), nil
}

var group singleflight.Group

func UsingSingleFlight(key string) {
	v, _, _ := group.Do(key, func() (interface{}, error) {
		return FetchExpensiveData()
	})
	fmt.Println(v)
}

func DoSingleFlight() {
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")

	time.Sleep(2 * time.Second)

	go UsingSingleFlight("key")
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")

	time.Sleep(2 * time.Second)
}

var (
	instance *Config
	once     sync.Once
	onceFunc = sync.OnceFunc(func() {
		fmt.Println("init config")
		instance = loadConfig()
	})
)

type Config struct{}

func loadConfig() *Config {
	return &Config{}
}
func GetConfig() *Config {
	// defer mu.Unlock()
	// mu.Lock()
	once.Do(func() {
		fmt.Println("init config")
		instance = loadConfig()
	})

	return instance
}

func GetConfigOnce() {
	onceFunc()
	onceFunc()
}
