package mocks

import (
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/mock"
)

var _ limiter.LimiterInterface = (*LimiterMock)(nil)

type LimiterMock struct {
	mock.Mock
}

func NewLimiterMock() *LimiterMock {
	return &LimiterMock{}
}

func (lm *LimiterMock) Limit(ip string, token string) error {
	args := lm.Called(ip, token)
	return args.Error(0)
}
