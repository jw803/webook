package ioc

import (
	"github.com/jw803/webook/internal/interface/event"
	"github.com/jw803/webook/internal/interface/event/article"
)

func NewConsumers(
	interactiveReadConsumer *article.InteractiveReadEventConsumer,
) map[string]event.Consumer {
	consumerMap := make(map[string]event.Consumer)
	consumerMap["interactive-read"] = interactiveReadConsumer
	return consumerMap
}
