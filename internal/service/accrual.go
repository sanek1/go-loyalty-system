package service

// package service

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"go.uber.org/zap"
// )

// // AccrualResponse представляет ответ от системы начисления баллов
// type AccrualResponse struct {
// 	Order   string  `json:"order"`
// 	Status  string  `json:"status"`
// 	Accrual float64 `json:"accrual,omitempty"`
// }

// // AccrualRepository интерфейс для работы с хранилищем начислений
// type AccrualRepository interface {
// 	SaveAccrual(ctx context.Context, orderNumber, status string, accrual float64) error
// }

// // AccrualClient интерфейс для работы с системой начисления баллов
// type AccrualClient interface {
// 	SendOrder(ctx context.Context, orderNumber string) error
// 	GetAccrual(ctx context.Context, orderNumber string) (*AccrualResponse, error)
// }

// // OrderAccrual сервис обработки начислений
// type OrderAccrual struct {
// 	ctx            context.Context
// 	logger         *zap.Logger
// 	repo           AccrualRepository
// 	client         AccrualClient
// 	processTimeout time.Duration
// 	retryInterval  time.Duration
// 	maxAttempts    int
// }

// // NewOrderAccrual создает новый сервис обработки начислений
// func NewOrderAccrual(ctx context.Context, logger *zap.Logger, repo AccrualRepository, client AccrualClient) *OrderAccrual {
// 	return &OrderAccrual{
// 		ctx:            ctx,
// 		logger:         logger,
// 		repo:           repo,
// 		client:         client,
// 		processTimeout: 30 * time.Second,
// 		retryInterval:  2 * time.Second,
// 		maxAttempts:    5,
// 	}
// }

// // processOrder обрабатывает заказ асинхронно
// func (op *OrderAccrual) processOrder(orderNumber string) {
// 	ctx, cancel := context.WithTimeout(op.ctx, op.processTimeout)
// 	defer cancel()

// 	op.logger.Info("processing order", zap.String("order", orderNumber))

// 	// Отправляем данные о заказе
// 	if err := op.sendOrderData(ctx, orderNumber); err != nil {
// 		op.handleProcessError(ctx, "send data", err, orderNumber)
// 		return
// 	}

// 	// Создаем каналы для результата и ошибок
// 	resultCh := make(chan *AccrualResponse, 1)
// 	errCh := make(chan error, 1)

// 	// Запускаем горутину для периодического опроса результата
// 	go func() {
// 		ticker := time.NewTicker(op.retryInterval)
// 		defer ticker.Stop()

// 		attempts := 0

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				errCh <- ctx.Err()
// 				return
// 			case <-ticker.C:
// 				attempts++

// 				accrualResp, err := op.getAccrualResult(ctx, orderNumber)
// 				if err != nil {
// 					if attempts >= op.maxAttempts {
// 						errCh <- fmt.Errorf("max attempts reached: %w", err)
// 						return
// 					}
// 					op.logger.Debug("retry getting accrual result",
// 						zap.String("order", orderNumber),
// 						zap.Int("attempt", attempts),
// 						zap.Error(err))
// 					continue
// 				}

// 				// Если статус PROCESSING, продолжаем опрос
// 				if accrualResp.Status == "PROCESSING" {
// 					if attempts >= op.maxAttempts {
// 						errCh <- fmt.Errorf("max attempts reached: still processing")
// 						return
// 					}
// 					op.logger.Debug("order still processing",
// 						zap.String("order", orderNumber),
// 						zap.Int("attempt", attempts))
// 					continue
// 				}

// 				// Если получили финальный статус, отправляем результат
// 				resultCh <- accrualResp
// 				return
// 			}
// 		}
// 	}()

// 	// Ожидаем результат или ошибку
// 	select {
// 	case accrualResp := <-resultCh:
// 		if err := op.repo.SaveAccrual(ctx, orderNumber, accrualResp.Status, accrualResp.Accrual); err != nil {
// 			op.handleProcessError(ctx, "save accrual", err, orderNumber)
// 			return
// 		}
// 		op.logger.Info("order processed successfully",
// 			zap.String("order", orderNumber),
// 			zap.String("status", accrualResp.Status),
// 			zap.Float64("accrual", accrualResp.Accrual))

// 	case err := <-errCh:
// 		op.handleProcessError(ctx, "process accrual", err, orderNumber)
// 		return

// 	case <-ctx.Done():
// 		op.handleProcessError(ctx, "context cancelled", ctx.Err(), orderNumber)
// 		return
// 	}
// }

// // sendOrderData отправляет данные о заказе в систему начисления
// func (op *OrderAccrual) sendOrderData(ctx context.Context, orderNumber string) error {
// 	return op.client.SendOrder(ctx, orderNumber)
// }

// // getAccrualResult получает результат обработки заказа
// func (op *OrderAccrual) getAccrualResult(ctx context.Context, orderNumber string) (*AccrualResponse, error) {
// 	return op.client.GetAccrual(ctx, orderNumber)
// }

// // handleProcessError обрабатывает ошибки процесса
// func (op *OrderAccrual) handleProcessError(ctx context.Context, msg string, err error, orderNumber string) {
// 	op.logger.Error(msg,
// 		zap.String("order", orderNumber),
// 		zap.Error(err))
// }
