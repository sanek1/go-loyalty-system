package app

import (
	"go-loyalty-system/internal/controller/accrual"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func NewPoolController(repo usecase.UserUseCase, address string, l *logging.ZapLogger) *accrual.OrderAccrual {
	orderProcessor := accrual.NewOrderProcessor(
		address,//"http://localhost:8081",
		3, // количество воркеров
		repo,
		l,
	)
	return orderProcessor
}

func startPool(orderProcessor *accrual.OrderAccrual) {
	// Запускаем процессор
	orderProcessor.Start()

	// При завершении приложения
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-gracefulShutdown
		orderProcessor.Stop()
	}()
}
