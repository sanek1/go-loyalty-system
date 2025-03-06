package services

// package services

// import (
// 	"context"
// 	"net/http"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/your-project/logger"
// 	"github.com/your-project/repositories"

// 	"github.com/your-project/utils/zap"
// )

// type OrderAccrual struct {
// 	ordersChan chan string
// 	wg         *sync.WaitGroup
// 	ctx        context.Context
// 	cancel     context.CancelFunc
// 	logger     *logger.Logger
// 	repo       repositories.AccrualRepository
// 	accrualURL string
// 	client     *http.Client
// 	semaphore  chan struct{}
// }

// func NewOrderProcessor(ctx context.Context, repo repositories.AccrualRepository, logger *logger.Logger,
// 	accrualURL string, workersCount int, maxParallelProcessing int) *OrderAccrual {
// 	ctx, cancel := context.WithCancel(ctx)

// 	// Если не указано, устанавливаем максимальное кол-во одновременных обработок как 10*workersCount
// 	if maxParallelProcessing <= 0 {
// 		maxParallelProcessing = 10 * workersCount
// 	}

// 	return &OrderAccrual{
// 		ctx:        ctx,
// 		cancel:     cancel,
// 		ordersChan: make(chan string, 100),
// 		repo:       repo,
// 		logger:     logger,
// 		accrualURL: accrualURL,
// 		client: &http.Client{
// 			Timeout: 10 * time.Second,
// 		},
// 		// Инициализируем семафор
// 		semaphore: make(chan struct{}, maxParallelProcessing),
// 	}
// }

// func (op *OrderAccrual) worker(id int) {
// 	defer op.wg.Done()
// 	op.logger.InfoCtx(context.Background(), "worker started", zap.Int("worker_id", id))

// 	for {
// 		select {
// 		case orderNumber, ok := <-op.ordersChan:
// 			if !ok {
// 				op.logger.InfoCtx(context.Background(), "worker stopped (channel closed)", zap.Int("worker_id", id))
// 				return
// 			}
// 			op.processOrder(orderNumber)
// 		case <-op.ctx.Done():
// 			op.logger.InfoCtx(context.Background(), "worker stopped (context canceled)", zap.Int("worker_id", id))
// 			return
// 		}
// 	}
// }

// func (op *OrderAccrual) processOrder(orderNumber string) {
// 	ctx, cancel := context.WithTimeout(op.ctx, processTimeout)
// 	defer cancel()
// 	op.logger.InfoCtx(ctx, "processing order", zap.String("order", orderNumber))

// 	// Сначала отправляем данные о заказе
// 	if err := op.sendOrderData(ctx, orderNumber); err != nil {
// 		op.handleProcessError(ctx, "send data ->"+err.Error(), err, orderNumber)
// 		return
// 	}

// 	// Запускаем асинхронную обработку для получения и сохранения результата
// 	go func() {
// 		// Получаем слот семафора
// 		select {
// 		case op.semaphore <- struct{}{}:
// 			// Получили слот, продолжаем
// 			defer func() {
// 				// Освобождаем слот по завершении
// 				<-op.semaphore
// 			}()
// 		case <-op.ctx.Done():
// 			// Контекст отменен, не запускаем обработку
// 			return
// 		}

// 		// Выполняем асинхронную обработку
// 		op.processOrderResult(orderNumber)
// 	}()

// 	op.logger.InfoCtx(ctx, "order sent for processing", zap.String("order", orderNumber))
// }

// // Новый метод для асинхронной обработки результатов
// func (op *OrderAccrual) processOrderResult(orderNumber string) {
// 	// Создаем новый контекст для асинхронной обработки
// 	ctx, cancel := context.WithTimeout(op.ctx, 30*time.Second)
// 	defer cancel()

// 	op.logger.InfoCtx(ctx, "async processing started", zap.String("order", orderNumber))

// 	// Получаем результат обработки с повторными попытками
// 	var accrualResp *AccrualResponse
// 	var err error

// 	// Добавляем механизм повторных попыток с экспоненциальной задержкой
// 	maxRetries := 5
// 	for i := 0; i < maxRetries; i++ {
// 		accrualResp, err = op.getAccrualResult(ctx, orderNumber)
// 		if err == nil {
// 			break
// 		}

// 		// Если статус PROCESSING, ждем и повторяем
// 		if strings.Contains(err.Error(), "processing") {
// 			waitTime := time.Duration(2<<i) * time.Second // Экспоненциальная задержка
// 			op.logger.InfoCtx(ctx, "order still processing, will retry",
// 				zap.String("order", orderNumber),
// 				zap.Duration("wait_time", waitTime),
// 				zap.Int("attempt", i+1))

// 			select {
// 			case <-time.After(waitTime):
// 				continue
// 			case <-ctx.Done():
// 				op.handleProcessError(ctx, "context canceled during retry", ctx.Err(), orderNumber)
// 				return
// 			}
// 		} else {
// 			// Другая ошибка, не связанная с обработкой
// 			op.handleProcessError(ctx, "get result failed", err, orderNumber)
// 			return
// 		}
// 	}

// 	if err != nil {
// 		op.handleProcessError(ctx, "get result failed after retries", err, orderNumber)
// 		return
// 	}

// 	// Сохраняем результат в базу
// 	if err := op.repo.SaveAccrual(ctx, orderNumber, accrualResp.Status, accrualResp.Accrual); err != nil {
// 		op.handleProcessError(ctx, "save accrual", err, orderNumber)
// 		return
// 	}

// 	op.logger.InfoCtx(ctx, "order processed successfully",
// 		zap.String("order", orderNumber),
// 		zap.String("status", accrualResp.Status),
// 		zap.Float64("accrual", accrualResp.Accrual))
// }
