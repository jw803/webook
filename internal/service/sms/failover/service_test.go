package failover

import (
	"context"
	"errors"
	"github.com/jw803/webook/internal/service/sms"
	smsmocks "github.com/jw803/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFailoverSMSService_Send(t *testing.T) {

	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		ctx     context.Context
		biz     string
		args    []string
		numbers []string

		wantErr error
	}{
		{
			name: "一次成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "biz1", []string{"123 "}, "09888xxxx").Return(nil)
				return []sms.Service{svc}
			},
			ctx:     context.Background(),
			biz:     "biz1",
			args:    []string{"123 "},
			numbers: []string{"09888xxxx"},
			wantErr: nil,
		},
		{
			name: "重試成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "biz1", []string{"123 "}, "09888xxxx").
					Return(errors.New("發送失敗"))
				svc2 := smsmocks.NewMockService(ctrl)
				svc2.EXPECT().Send(gomock.Any(), "biz1", []string{"123 "}, "09888xxxx").Return(nil)
				return []sms.Service{svc, svc2}
			},
			ctx:     context.Background(),
			biz:     "biz1",
			args:    []string{"123 "},
			numbers: []string{"09888xxxx"},
			wantErr: nil,
		},
		{
			name: "重試全部失敗",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "biz1", []string{"123 "}, "09888xxxx").
					Return(errors.New("發送失敗"))
				svc2 := smsmocks.NewMockService(ctrl)
				svc2.EXPECT().Send(gomock.Any(), "biz1", []string{"123 "}, "09888xxxx").
					Return(errors.New("發送還是失敗"))
				return []sms.Service{svc, svc2}
			},
			ctx:     context.Background(),
			biz:     "biz1",
			args:    []string{"123 "},
			numbers: []string{"09888xxxx"},
			wantErr: errors.New("全部服务商都失败了"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svcs := tc.mock(ctrl)
			failoverSmsService := NewFailoverSMSService(svcs)
			err := failoverSmsService.Send(tc.ctx, tc.biz, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
