package user

import (
	"github.com/jw803/webook/pkg/loggerx"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestNewUserHandler(t *testing.T) {
	suite.Run(t, new(userHandlerSuite))
}

type userHandlerSuite struct {
	suite.Suite

	logger loggerx.Logger
}

func (s *userHandlerSuite) SetupSuite() {
	s.logger = loggerx.NewLocalLogger()
}
