package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-loyalty-system/internal/entity"
	"go-loyalty-system/internal/usecase"
	"go-loyalty-system/pkg/logging"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type OrderAccrual struct {
	client     *http.Client
	baseURL    string
	numWorkers int
	ordersChan chan string
	logger     *logging.ZapLogger
	repo       usecase.UserUseCase
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	running    bool
	mu         sync.Mutex
}

const (
	defaultBufferSize = 100
	defaultTimeout    = 10 * time.Second
	processTimeout    = 60 * time.Second
	collectorInterval = 20 * time.Second
	maxRetries        = 5
	initialBackoff    = time.Second
	maxBackoff        = 30 * time.Second
)

func NewOrderProcessor(baseURL string, numWorkers int, repo usecase.UserUseCase, l *logging.ZapLogger) *OrderAccrual {
	ctx, cancel := context.WithCancel(context.Background())
	return &OrderAccrual{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL:    baseURL,
		numWorkers: numWorkers,
		ordersChan: make(chan string, defaultBufferSize),
		logger:     l,
		repo:       repo,
		ctx:        ctx,
		cancel:     cancel,
		running:    false,
	}
}

func (op *OrderAccrual) Start() {
	op.logger.InfoCtx(op.ctx, "service started")
	op.logger.InfoCtx(op.ctx, "baseURL -> "+op.baseURL)
	op.mu.Lock()
	defer op.mu.Unlock()

	if op.running {
		return
	}

	op.running = true
	op.startWorkers()
	op.startCollector()
}
func (op *OrderAccrual) startWorkers() {
	for i := 0; i < op.numWorkers; i++ {
		op.wg.Add(1)
		go op.worker(i)
	}
	op.logger.InfoCtx(op.ctx, "workers started", zap.Int("count", op.numWorkers))
}
func (op *OrderAccrual) startCollector() {
	op.wg.Add(1)
	go op.collectUnprocessedOrders()
	op.logger.InfoCtx(op.ctx, "collector started")
}

func (op *OrderAccrual) Stop() {
	op.mu.Lock()
	defer op.mu.Unlock()

	if !op.running {
		return
	}

	op.cancel()
	op.wg.Wait()
	op.running = false
	op.logger.InfoCtx(op.ctx, "service stopped")
}

func (op *OrderAccrual) AddOrder(orderNumber string) {
	select {
	case op.ordersChan <- orderNumber:
		op.logger.InfoCtx(context.Background(), "order added to processing queue #"+orderNumber, zap.String("order", orderNumber))
	default:
		// Если канал переполнен, логируем это
		op.logger.WarnCtx(context.Background(), "orders channel is full, order skipped", zap.String("order", orderNumber))
		// Запускаем горутину для повторной попытки
		go op.retryAddOrder(orderNumber)
	}
}

func (op *OrderAccrual) retryAddOrder(orderNumber string) {
	backoff := time.Second
	maxBackoff := time.Second * 30

	for {
		select {
		case <-op.ctx.Done():
			return
		case <-time.After(backoff):
			select {
			case op.ordersChan <- orderNumber:
				op.logger.InfoCtx(context.Background(), "order added to processing queue after retry",
					zap.String("order", orderNumber))
				return
			default:
				// Увеличиваем время ожидания
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				continue
			}
		}
	}
}

func (op *OrderAccrual) worker(id int) {
	defer op.wg.Done()
	op.logger.InfoCtx(context.Background(), "worker started", zap.Int("worker_id", id))

	for {
		select {
		case orderNumber, ok := <-op.ordersChan:
			if !ok {
				op.logger.InfoCtx(context.Background(), "worker stopped (channel closed)", zap.Int("worker_id", id))
				return
			}
			op.processOrder(orderNumber)
		case <-op.ctx.Done():
			op.logger.InfoCtx(context.Background(), "worker stopped (context canceled)", zap.Int("worker_id", id))
			return
		}
	}
}

func (op *OrderAccrual) handleProcessError(ctx context.Context, stage string, err error, orderNumber string) {
	op.logger.ErrorCtx(ctx, "processing failed "+stage,
		zap.String("stage", stage),
		zap.Error(err),
		zap.String("order", orderNumber))
	go op.retryOrder(orderNumber)
}

func (op *OrderAccrual) processOrder(orderNumber string) {
	ctx, cancel := context.WithTimeout(op.ctx, processTimeout)
	defer cancel()
	op.logger.InfoCtx(ctx, "processing order", zap.String("order", orderNumber))

	// Сначала отправляем данные о заказе
	if err := op.sendOrderData(ctx, orderNumber); err != nil {
		op.handleProcessError(ctx, "send data ->"+err.Error(), err, orderNumber)
		return
	}

	// Запускаем асинхронную обработку для получения и сохранения результата
	go op.processOrderResult(orderNumber)

	op.logger.InfoCtx(ctx, "order sent for processing", zap.String("order", orderNumber))
}

// Новый метод для асинхронной обработки результатов
func (op *OrderAccrual) processOrderResult(orderNumber string) {
	// Создаем новый контекст для асинхронной обработки
	ctx, cancel := context.WithTimeout(op.ctx, maxBackoff)
	defer cancel()

	op.logger.InfoCtx(ctx, "async processing started", zap.String("order", orderNumber))

	// Получаем результат обработки с повторными попытками
	var accrualResp *entity.AccrualResponse
	var err error

	// Добавляем механизм повторных попыток с экспоненциальной задержкой
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		accrualResp, err = op.getAccrualResult(ctx, orderNumber)
		if err == nil {
			break
		}

		// Если статус PROCESSING, ждем и повторяем
		if strings.Contains(err.Error(), "processing") {
			waitTime := time.Duration(2<<i) * time.Second // Экспоненциальная задержка
			op.logger.InfoCtx(ctx, "order still processing, will retry",
				zap.String("order", orderNumber),
				zap.Duration("wait_time", waitTime),
				zap.Int("attempt", i+1))

			select {
			case <-time.After(waitTime):
				continue
			case <-ctx.Done():
				op.handleProcessError(ctx, "context canceled during retry", ctx.Err(), orderNumber)
				return
			}
		} else {
			op.handleProcessError(ctx, "get result failed", err, orderNumber)
			return
		}
	}

	if err != nil {
		op.handleProcessError(ctx, "get result failed after retries", err, orderNumber)
		return
	}

	// Сохраняем результат в базу
	if err := op.repo.SaveAccrual(ctx, orderNumber, accrualResp.Status, accrualResp.Accrual); err != nil {
		op.handleProcessError(ctx, "save accrual", err, orderNumber)
		return
	}

	op.logger.InfoCtx(ctx, "order processed successfully",
		zap.String("order", orderNumber),
		zap.String("status", accrualResp.Status))
}

func (op *OrderAccrual) checkResponse(resp *http.Response) error {
	ctx := context.Background()
	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		op.logger.InfoCtx(ctx, "order processed successfully")
		return nil
	case http.StatusTooManyRequests:
		return op.logAndReturnError(ctx, "rate limit exceeded", nil)
	case http.StatusBadRequest:
		return fmt.Errorf("bad request")
	case http.StatusInternalServerError:
		return fmt.Errorf("server error")
	case http.StatusNoContent:
		return op.logAndReturnError(ctx, "order not found", nil)

	default:
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}

func (op *OrderAccrual) getAccrualResult(ctx context.Context, orderNumber string) (*entity.AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", op.baseURL, orderNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, op.logAndReturnError(ctx, "failed to create request", err)
	}

	resp, err := op.client.Do(req)
	if err != nil {
		return nil, op.logAndReturnError(ctx, "failed to get order info", err)
	}
	defer resp.Body.Close()

	if err := op.checkResponse(resp); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, op.logAndReturnError(ctx, "unexpected status code ->"+resp.Status, nil)
	}

	var accrualResp entity.AccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&accrualResp); err != nil {
		return nil, op.logAndReturnError(ctx, "failed to decode response", err)
	}
	return &accrualResp, nil
}

func (op *OrderAccrual) createRequest(ctx context.Context, method, path string, body []byte) (*http.Request, error) {
	op.logger.InfoCtx(ctx, "path->"+path)
	op.logger.InfoCtx(ctx, "creating request"+op.baseURL+path)
	req, err := http.NewRequestWithContext(ctx, method, op.baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		op.logger.ErrorCtx(ctx, "create request error", zap.Error(err))
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (op *OrderAccrual) sendOrderData(ctx context.Context, orderNumber string) error {
	orderData := entity.AccrualOrder{Order: orderNumber, Goods: []entity.Product{
		{
			Description: "test",
			Price:       100,
		},
		{
			Description: "test2",
			Price:       200,
		},
	}}
	jsonData, err := json.Marshal(orderData)
	if err != nil {
		op.logger.ErrorCtx(ctx, "marshal error", zap.Error(err))
		return err
	}
	op.logger.InfoCtx(ctx, "sending order data "+orderNumber, zap.String("order", orderNumber))
	op.logger.InfoCtx(ctx, "order data", zap.String("order", string(jsonData)))
	req, err := op.createRequest(ctx, "POST", "/api/orders", jsonData)
	if err != nil {
		op.logger.ErrorCtx(ctx, "create request error ->"+err.Error(), zap.Error(err))
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := op.client.Do(req)
	if err != nil {
		return op.logAndReturnError(ctx, "sendOrderData - failed to send request", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return op.logAndReturnError(ctx, "sendOrderData - rate limit exceeded", nil)
	case http.StatusNoContent:
		op.logger.InfoCtx(ctx, "order already registered")
		return nil
	case http.StatusAccepted:
		op.logger.InfoCtx(ctx, "The order has been successfully accepted for processing")
		return nil
	case http.StatusConflict:
		op.logger.InfoCtx(ctx, "The order is already processed")
		return nil
	case http.StatusBadRequest:
		return op.logAndReturnError(ctx, "sendOrderData - bad request", nil)
	case http.StatusInternalServerError:
		return op.logAndReturnError(ctx, "sendOrderData - internal server error", nil)
	}
	if resp.StatusCode != http.StatusOK {
		return op.logAndReturnError(ctx, "sendOrderData - unexpected status code -"+resp.Status, nil)
	}
	return nil
}

func (op *OrderAccrual) retryOrder(orderNumber string) {
	backoff := time.Second
	maxBackoff := maxBackoff
	maxRetries := maxRetries
	retries := 0

	for retries < maxRetries {
		select {
		case <-op.ctx.Done():
			return
		case <-time.After(backoff):
			op.AddOrder(orderNumber)
			retries++
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

func (op *OrderAccrual) collectUnprocessedOrders() {
	ctx := context.Background()
	defer op.wg.Done()
	ticker := time.NewTicker(collectorInterval)
	defer ticker.Stop()
	op.logger.InfoCtx(ctx, ("unprocessed orders collector started"))

	for {
		select {
		case <-op.ctx.Done():
			op.logger.InfoCtx(ctx, "unprocessed orders collector stopped")
			return
		case <-ticker.C:

			op.checkUnprocessedOrders()
			op.logger.InfoCtx(ctx, "⌛ No orders found...")
		}
	}
}

func (op *OrderAccrual) checkUnprocessedOrders() {
	orders, err := op.repo.GetUnprocessedOrders(op.ctx)
	if err != nil {
		op.logger.ErrorCtx(op.ctx, "failed to get unprocessed orders", zap.Error(err))
		return
	}

	if len(orders) > 0 {
		op.logger.InfoCtx(op.ctx, "processing unprocessed orders", zap.Int("count", len(orders)))
		for _, order := range orders {
			op.AddOrder(order)
		}
	}
}

func (op *OrderAccrual) logAndReturnError(ctx context.Context, method string, err error) error {
	msg := fmt.Sprintf("%s - %s: %v", "Accrual", method, err)
	op.logger.ErrorCtx(ctx, msg, zap.Error(err))
	return err
}
