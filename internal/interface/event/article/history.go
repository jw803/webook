package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jw803/webook/internal/domain"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/jw803/webook/pkg/samarax"
	"time"
)

type HistoryRecordConsumer struct {
	repo   repository.HistoryRecordRepository
	client sarama.Client
	l      loggerx.Logger
}

func (i *HistoryRecordConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		ctx := context.Background()
		err := cg.Consume(ctx,
			[]string{TopicReadEvent},
			samarax.NewHandler[ReadEvent](i.l, i.Consume))
		if err != nil {
			i.l.Error(ctx, "退出消费", loggerx.Error(err))
		}
	}()
	return err
}

func (i *HistoryRecordConsumer) Consume(msg *sarama.ConsumerMessage, event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.AddRecord(ctx, domain.HistoryRecord{
		BizId: event.Aid,
		Biz:   "article",
		Uid:   event.Uid,
	})
}
