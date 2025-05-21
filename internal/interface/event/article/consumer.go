package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/jw803/webook/pkg/samarax"
	"time"
)

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
	l      loggerx.Logger
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository,
	client sarama.Client, l loggerx.Logger) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

func (i *InteractiveReadEventConsumer) Start() error {
	ctx := context.Background()
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}

	er := cg.Consume(context.Background(),
		[]string{TopicReadEvent},
		samarax.NewBatchHandler[ReadEvent](i.l, i.BatchConsume))
	if er != nil {
		i.l.Error(ctx, "退出消费", loggerx.Error(er))
	}

	return err
}

func (i *InteractiveReadEventConsumer) StartV1() error {
	ctx := context.Background()
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(ctx,
			[]string{TopicReadEvent},
			samarax.NewHandler[ReadEvent](i.l, i.Consume))
		if er != nil {
			i.l.Error(ctx, "退出消费", loggerx.Error(er))
		}
	}()
	return err
}
func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	events []ReadEvent) error {
	bizs := make([]string, 0, len(events))
	bizIds := make([]int64, 0, len(events))
	for _, evt := range events {
		bizs = append(bizs, "article")
		bizIds = append(bizIds, evt.Aid)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.BatchIncrReadCnt(ctx, bizs, bizIds)
}

func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage,
	event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.IncrReadCnt(ctx, "article", event.Aid)
}
