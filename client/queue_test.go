package client

//go:generate mockery --output . --filename  mock_service_test.go --outpkg client --srcpkg github.com/kozmod/omni/external --name Service  --structname mockService

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/kozmod/omni/external"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

const (
	externalServiceGetLimitsFunc = "GetLimits"
	externalServiceProcessFunc   = "Process"
)

func Test_(t *testing.T) {
	var (
		item1 = Item{
			ID: uuid.NewV4(),
		}
		item2 = Item{
			ID: uuid.NewV4(),
		}
		item3 = Item{
			ID: uuid.NewV4(),
		}
	)
	t.Run("process", func(t *testing.T) {
		t.Run("in_the_less_then_batch", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {

				const (
					batchSize = uint64(10)
					delay     = 10 * time.Second
				)

				t.Run("single", func(t *testing.T) {
					var (
						ctx         = context.Background()
						serviceMock = new(mockService)
						q, _        = NewExternalServerQueue(serviceMock, 0)

						expWait = 1 * time.Second
					)

					serviceMock.On(externalServiceGetLimitsFunc).
						Return(batchSize, delay).
						Times(1)
					serviceMock.On(externalServiceProcessFunc, ctx, external.Batch{struct{}{}}).
						Return(nil).
						Times(1)

					q.AddProcess(item1)

					processed, wait, err := q.process(ctx)
					assert.Nil(t, err)
					assert.Equal(t, expWait, wait)
					assert.Equal(t, []ID{item1.ID}, processed)

					serviceMock.AssertExpectations(t)
				})
				t.Run("multi", func(t *testing.T) {
					var (
						ctx         = context.Background()
						serviceMock = new(mockService)
						q, _        = NewExternalServerQueue(serviceMock, 0)

						insetItems = []Item{item1, item2, item3}
						expWait    = time.Duration(len(insetItems)) * time.Second
					)

					serviceMock.On(externalServiceGetLimitsFunc).
						Return(batchSize, delay).
						Times(1)
					serviceMock.On(externalServiceProcessFunc,
						ctx,
						external.Batch{
							struct{}{},
							struct{}{},
							struct{}{},
						}).
						Return(nil).
						Times(1)

					q.AddProcess(insetItems...)

					processed, wait, err := q.process(ctx)
					assert.Nil(t, err)
					assert.Equal(t, expWait, wait)
					assert.Equal(t,
						[]ID{
							item1.ID,
							item2.ID,
							item3.ID,
						},
						processed)

					serviceMock.AssertExpectations(t)
				})
			})

			t.Run("error", func(t *testing.T) {

				var (
					expWait = 0 * time.Second
				)

				t.Run("when_limit_is_zero", func(t *testing.T) {
					var (
						ctx         = context.Background()
						serviceMock = new(mockService)
						q, _        = NewExternalServerQueue(serviceMock, 0)
					)

					serviceMock.On(externalServiceGetLimitsFunc).
						Return(uint64(0), 100*time.Second).
						Times(1)

					q.AddProcess(item1)

					processed, wait, err := q.process(ctx)
					assert.ErrorIs(t, err, ErrExternalServiceZeroLimit)
					assert.Equal(t, expWait, wait)
					assert.Empty(t, processed)

					serviceMock.AssertExpectations(t)
				})
				t.Run("when_delay_is_zero", func(t *testing.T) {
					var (
						ctx         = context.Background()
						serviceMock = new(mockService)
						q, _        = NewExternalServerQueue(serviceMock, 0)
					)

					serviceMock.On(externalServiceGetLimitsFunc).
						Return(uint64(100), time.Duration(0)).
						Times(1)

					q.AddProcess(item1)

					processed, wait, err := q.process(ctx)
					assert.ErrorIs(t, err, ErrExternalServiceZeroDelay)
					assert.Equal(t, expWait, wait)
					assert.Empty(t, processed)

					serviceMock.AssertExpectations(t)
				})
				t.Run("when_external_server_return_error", func(t *testing.T) {
					var (
						ctx         = context.Background()
						serviceMock = new(mockService)
						q, _        = NewExternalServerQueue(serviceMock, 0)

						expErr = fmt.Errorf("SOME_external_server_err")
					)

					serviceMock.On(externalServiceGetLimitsFunc).
						Return(uint64(10), 10*time.Second).
						Times(1)
					serviceMock.On(externalServiceProcessFunc, ctx, mock.Anything).
						Return(expErr).
						Times(1)

					q.AddProcess(item1)

					processed, wait, err := q.process(ctx)
					assert.Equal(t, expErr, err)
					assert.Equal(t, expWait, wait)
					assert.Empty(t, processed)

					serviceMock.AssertExpectations(t)
				})
			})
		})
	})
}
