package ioc

import (
	"github.com/jw803/webook/internal/interface/event"
)

func NewConsumers() map[string]event.Consumer {
	consumerMap := make(map[string]event.Consumer)
	return consumerMap
}
