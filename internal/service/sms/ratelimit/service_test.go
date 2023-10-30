package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"github.com/jw803/webook/internal/service/sms"
	smsmocks "github.com/jw803/webook/internal/service/sms/mocks"
	"github.com/jw803/webook/pkg/ratelimit"
	limitmocks "github.com/jw803/webook/pkg/ratelimit/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRatelimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter)
		ctx     context.Context
		tpl     string
		args    []string
		numbers []string
		wantErr error
	}{
		{
			name: "正常發送",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := limitmocks.NewMockLimiter(ctrl)

				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().
					Send(gomock.Any(), "string", []string{"string1", "string2"}, "string1").
					Return(nil)
				return svc, limiter
			},
			ctx:     context.Background(),
			tpl:     "string",
			args:    []string{"string1", "string2"},
			numbers: []string{"string1"},
			wantErr: nil,
		},
		{
			name: "觸發限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return nil, limiter
			},
			ctx:     context.Background(),
			tpl:     "string",
			args:    []string{"string1", "string2"},
			numbers: []string{"string1"},
			wantErr: errLimited,
		},
		{
			name: "限流器異常",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, errors.New("限流器異常"))
				return nil, limiter
			},
			ctx:     context.Background(),
			tpl:     "string",
			args:    []string{"string1", "string2"},
			numbers: []string{"string1"},
			wantErr: fmt.Errorf("短信服务判断是否限流出现问题，%w", errors.New("限流器異常")),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc, limiter := tc.mock(ctrl)

			s := &RatelimitSMSService{
				svc:     svc,
				limiter: limiter,
			}

			err := s.Send(tc.ctx, tc.tpl, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
