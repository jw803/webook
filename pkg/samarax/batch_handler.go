package samarax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/jw803/webook/pkg/loggerx"
	"time"
)

type BatchHandler[T any] struct {
	SamaraxBaseHandler
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	l  loggerx.Logger
}

func NewBatchHandler[T any](l loggerx.Logger, fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{fn: fn, l: l}
}

func (h *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx, err := h.GetContext()
	if err != nil {
		h.l.P1(ctx, "failed to generate ctx")
		return err
	}

	msgs := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		var done = false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-timeoutCtx.Done():
				// 超时了
				done = true
			case msg, ok := <-msgs:
				if !ok {
					cancel()
					return nil
				}
				batch = append(batch, msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					h.l.Error(ctx, "failed to serialize the event msg",
						loggerx.String("topic", msg.Topic),
						loggerx.Int32("partition", msg.Partition),
						loggerx.Int64("offset", msg.Offset),
						loggerx.Error(err))
					continue
				}
				batch = append(batch, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		// 凑够了一批，然后你就处理
		err := h.fn(batch, ts)
		if err != nil {
			h.l.Error(ctx, "failed to process the event msg", loggerx.Error(err))
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
