package client

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kozmod/omni/external"
)

type ExternalServerQueue struct {
	service external.Service

	errChCap int
	errCh    chan error

	mx    sync.Mutex
	once  sync.Once
	queue *list.List
}

func NewExternalServerQueue(service external.Service, errChCap int) (*ExternalServerQueue, <-chan error) {
	if errChCap < 0 {
		errChCap = 0
	}
	q := ExternalServerQueue{
		service:  service,
		queue:    list.New(),
		errChCap: errChCap,
	}
	return &q, q.errCh
}

func (c *ExternalServerQueue) AddProcess(items ...Item) {
	if len(items) == 0 {
		return
	}
	c.mx.Lock()
	defer c.mx.Unlock()

	for _, item := range items {
		c.queue.PushBack(item)
	}

}

func (c *ExternalServerQueue) Start(ctx context.Context) {
	c.once.Do(func() {
		c.errCh = make(chan error, c.errChCap)
		c.startProcess(ctx)
		close(c.errCh)
	})
}

func (c *ExternalServerQueue) startProcess(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				processed, wait, err := c.process(ctx)
				if err != nil {
					c.errCh <- fmt.Errorf("error accured on items [%v]: %w", processed, err)
					return
				}
				<-time.NewTicker(wait).C
			}
		}
	}()
}

func (c *ExternalServerQueue) process(ctx context.Context) ([]ID, time.Duration, error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	limit, delay := c.service.GetLimits()
	switch {
	case limit == 0:
		return nil, 0, ErrExternalServiceZeroLimit
	case delay == 0:
		return nil, 0, ErrExternalServiceZeroDelay
	}

	processOneItemDelay := uint64(delay) / limit

	var (
		batch      = make(external.Batch, 0, int(limit))
		itemsUUIDs = make([]ID, 0, int(limit))
	)

	for element, i := c.queue.Front(), uint64(0); i <= limit && element != nil; element, i = c.queue.Front(), i+1 {
		item := c.queue.Remove(element).(Item)
		batch = append(batch, mapItemToDTO(item))
		itemsUUIDs = append(itemsUUIDs, item.ID)
	}

	if len(batch) == 0 {
		return nil, 0, nil
	}

	err := c.service.Process(ctx, batch)
	switch {
	case errors.Is(err, external.ErrBlocked):
		return nil, 0, ErrClientWasBlocked
	case err != nil:
		return nil, 0, err
	default:
		return itemsUUIDs, time.Duration(len(batch)) * time.Duration(processOneItemDelay), nil
	}
}

func mapItemToDTO(_ Item) external.Item {
	// required mapping for domain entity to DTO
	return external.Item{}
}
