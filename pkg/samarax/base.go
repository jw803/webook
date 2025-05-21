package samarax

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jw803/webook/pkg/trace_id"
	"github.com/jw803/webook/pkg/uuidx"
)

type SamaraxBaseHandler struct {
	uuidFunc uuidx.UuidFn
}

func NewSamaraxBaseHandler(uuidFunc uuidx.UuidFn) *SamaraxBaseHandler {
	return &SamaraxBaseHandler{uuidFunc: uuidFunc}
}

func (h *SamaraxBaseHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *SamaraxBaseHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *SamaraxBaseHandler) GetContext() (context.Context, error) {
	ctx := context.Background()
	uuid, err := h.uuidFunc()
	if err != nil {
		return ctx, err
	}
	ctxWithTraceID := context.WithValue(context.Background(), trace_id.TraceIDKey, uuid)
	return ctxWithTraceID, nil
}
