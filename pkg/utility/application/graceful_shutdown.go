package application

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	completed = "Completed graceful shutdown of the app"
	reason    = "Received signal: %v"
)

var (
	shutDownOnce sync.Once
	terminate    = make(chan os.Signal, 1)
	waitGroup    sync.WaitGroup
)

func GracefulShutdown() {
	shutDownOnce.Do(func() {
		signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
		waitGroup.Add(1)

		go func() {
			code := <-terminate
			logging.Logger(fmt.Sprintf(reason, code.String()))
			os.Exit(1)
		}()
	})
	logging.Logger(completed)
	waitGroup.Wait()
}
