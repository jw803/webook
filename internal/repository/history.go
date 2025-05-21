package repository

import (
	"context"
	"github.com/jw803/webook/internal/domain"
)

type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
