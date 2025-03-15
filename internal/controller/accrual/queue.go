package accrual

import (
	"container/list"
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// OrderQueue представляет in-memory очередь заказов
type OrderQueue struct {
	items    *list.List
	capacity int
	mu       sync.RWMutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	logger   *zap.Logger
}

// QueueItem представляет элемент очереди
type QueueItem struct {
	OrderNumber string
	Attempts    int
	LastAttempt time.Time
}

// NewOrderQueue создает новую очередь заказов
func NewOrderQueue(capacity int, logger *zap.Logger) *OrderQueue {
	q := &OrderQueue{
		items:    list.New(),
		capacity: capacity,
		logger:   logger,
	}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

// Enqueue добавляет заказ в очередь
func (q *OrderQueue) Enqueue(ctx context.Context, orderNumber string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Проверяем, есть ли уже такой заказ в очереди
	for e := q.items.Front(); e != nil; e = e.Next() {
		if item := e.Value.(*QueueItem); item.OrderNumber == orderNumber {
			q.logger.Debug("order already in queue",
				zap.String("order", orderNumber),
				zap.Int("attempts", item.Attempts))
			return nil
		}
	}

	// Ждем, пока появится место в очереди
	for q.items.Len() >= q.capacity {
		done := make(chan struct{})
		go func() {
			q.notFull.Wait()
			close(done)
		}()

		q.mu.Unlock()
		select {
		case <-ctx.Done():
			q.mu.Lock()
			return ctx.Err()
		case <-done:
			q.mu.Lock()
		}
	}

	item := &QueueItem{
		OrderNumber: orderNumber,
		Attempts:    0,
		LastAttempt: time.Now(),
	}
	q.items.PushBack(item)
	q.logger.Debug("order enqueued",
		zap.String("order", orderNumber),
		zap.Int("queue_size", q.items.Len()))

	// Сигнализируем, что очередь не пуста
	q.notEmpty.Signal()
	return nil
}

// Dequeue извлекает заказ из очереди
func (q *OrderQueue) Dequeue(ctx context.Context) (*QueueItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Ждем, пока в очереди появятся элементы
	for q.items.Len() == 0 {
		done := make(chan struct{})
		go func() {
			q.notEmpty.Wait()
			close(done)
		}()

		q.mu.Unlock()
		select {
		case <-ctx.Done():
			q.mu.Lock()
			return nil, ctx.Err()
		case <-done:
			q.mu.Lock()
		}
	}

	element := q.items.Front()
	item := element.Value.(*QueueItem)
	q.items.Remove(element)

	q.logger.Debug("order dequeued",
		zap.String("order", item.OrderNumber),
		zap.Int("queue_size", q.items.Len()))

	// Сигнализируем, что в очереди появилось место
	q.notFull.Signal()
	return item, nil
}

// RetryOrder помещает заказ обратно в очередь с увеличенным счетчиком попыток
func (q *OrderQueue) RetryOrder(ctx context.Context, item *QueueItem) error {
	item.Attempts++
	item.LastAttempt = time.Now()
	return q.Enqueue(ctx, item.OrderNumber)
}

// Size возвращает текущий размер очереди
func (q *OrderQueue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.items.Len()
}

// Clear очищает очередь
func (q *OrderQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items.Init()
	q.notFull.Broadcast()
}
