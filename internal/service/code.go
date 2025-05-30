package service

import (
	"context"
	"fmt"
	"github.com/jw803/webook/internal/repository"
	"github.com/jw803/webook/internal/service/sms"
	"go.uber.org/atomic"
	"math/rand"
)

var codeTplId = atomic.String{}

type CodeService interface {
	Send(ctx context.Context,
		// 区别业务场景
		biz string, phone string) error
	Verify(ctx context.Context, biz string,
		phone string, inputCode string) (bool, error)
}

type codeService struct {
	codeRepo repository.CodeRepository
	smsSvc   sms.Service
	//tplId string
}

func NewSMSCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	codeTplId.Store("1877556")
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	codeTplId.Store(viper.GetString("code.tpl.id"))
	//})

	return &codeService{
		codeRepo: repo,
		smsSvc:   smsSvc,
	}
}

// Send 发验证码，我需要什么参数？
func (svc *codeService) Send(ctx context.Context,
	// 区别业务场景
	biz string,
	phone string) error {
	// 生成一个验证码
	code := svc.generateCode()
	// 塞进去 Redis
	err := svc.codeRepo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	err = svc.smsSvc.Send(ctx, codeTplId.Load(), []string{code}, phone)
	if err != nil {
		return err
	}
	//if err != nil {
	// 这个地方怎么办？
	// 这意味着，Redis 有这个验证码，但是不好意思，
	// 我能不能删掉这个验证码？
	// 你这个 err 可能是超时的 err，你都不知道，发出了没
	// 在这里重试
	// 要重试的话，初始化的时候，传入一个自己就会重试的 smsSvc
	//}
	return nil
}

func (svc *codeService) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {
	return svc.codeRepo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%06d", num)
}

//func (svc *codeService) VerifyV1(ctx context.Context, biz string,
//	phone string, inputCode string) error {
//
//}
