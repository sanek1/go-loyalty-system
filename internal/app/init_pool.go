package app

import (
	"go-loyalty-system/internal/controller/accrual"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

const (
	numWorkers = 5
)

func NewPoolController(repo usecase.UserUseCase, address string, l *logging.ZapLogger) *accrual.OrderAccrual {
	orderProcessor := accrual.NewOrderProcessor(
		address,
		numWorkers,
		repo,
		l,
	)
	return orderProcessor
}

func startPool(orderProcessor *accrual.OrderAccrual) {
	orderProcessor.Start()
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-gracefulShutdown
		orderProcessor.Stop()
	}()
}
