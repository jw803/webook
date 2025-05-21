package samarax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/jw803/webook/pkg/loggerx"
)

type Handler[T any] struct {
	SamaraxBaseHandler
	l  loggerx.Logger
	fn func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](l loggerx.Logger, fn func(msg *sarama.ConsumerMessage, event T) error) *Handler[T] {
	return &Handler[T]{l: l, fn: fn}
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx, err := h.GetContext()
	if err != nil {
		h.l.P1(ctx, "failed to generate ctx")
		return err
	}

	msgs := claim.Messages()
	for msg := range msgs {
		// 在这里调用业务处理逻辑
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			// 你也可以在这里引入重试的逻辑
			h.l.Error(ctx, "failed to serialize the event msg",
				loggerx.String("topic", msg.Topic),
				loggerx.Int32("partition", msg.Partition),
				loggerx.Int64("offset", msg.Offset),
				loggerx.Error(err))
		}
		err = h.fn(msg, t)
		if err != nil {
			h.l.Error(ctx, "failed to process the event msg",
				loggerx.String("topic", msg.Topic),
				loggerx.Int32("partition", msg.Partition),
				loggerx.Int64("offset", msg.Offset),
				loggerx.Error(err))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
