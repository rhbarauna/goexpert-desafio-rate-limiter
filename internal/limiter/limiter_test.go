package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/configs"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/limiter"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/mocks"
)

type LimiterSuite struct {
	suite.Suite
	limiter     limiter.LimiterInterface
	storage     *mocks.StorageMock
	cooldown    int
	maxRequests int
	ttl         int
	tokens      map[string]configs.TokenConfig
}

func (suite *LimiterSuite) SetupTest() {
	// Configurar a instância do RateLimiter para cada teste
	suite.storage = mocks.NewStorageMock()
	suite.tokens = map[string]configs.TokenConfig{
		"default": {
			Name:        "default",
			MaxRequests: 2,
			Cooldown:    5,
		},
		"custom": {
			Name:        "custom",
			MaxRequests: 3,
			Cooldown:    8,
		},
	}
	suite.cooldown = 5
	suite.maxRequests = 2
	suite.ttl = 10

	suite.limiter = limiter.NewLimiter(suite.storage, suite.cooldown, suite.maxRequests, suite.ttl, suite.tokens)
}

func (suite *LimiterSuite) TestLimit_CustomToken_Success() {
	ctx := context.Background()
	suite.storage.On("IsBlocked", ctx, "ratelimit:blocked:custom").Return(false, nil)
	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:custom", suite.ttl).Return(1, nil)
	err := suite.limiter.Limit("127.0.0.1", "custom")
	assert.NoError(suite.T(), err)
}

func (suite *LimiterSuite) TestLimit_ExceedTokenMaxRequests_Error() {
	tknCfg := suite.tokens["default"]
	ctx := context.Background()

	suite.storage.On("IsBlocked", ctx, "ratelimit:blocked:default").Return(false, nil)

	for i := 0; i < tknCfg.MaxRequests; i++ {
		suite.storage.On("Increment", ctx, "ratelimit:req_qnt:default", suite.ttl).Return(i+1, nil).Once()
		err := suite.limiter.Limit("127.0.0.1", tknCfg.Name)
		assert.NoError(suite.T(), err)
	}

	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:default", suite.ttl).Return(tknCfg.MaxRequests+1, nil).Once()
	suite.storage.On("Set", ctx, "ratelimit:blocked:default", tknCfg.Cooldown).Return(nil)
	err := suite.limiter.Limit("127.0.0.1", tknCfg.Name)
	assert.EqualError(suite.T(), err, limiter.ErrLimitedAccess.Error())
}

func (suite *LimiterSuite) TestLimit_TokenBlockingDuration() {
	tknCfg := suite.tokens["default"]
	ctx := context.Background()

	suite.storage.On("IsBlocked", ctx, "ratelimit:blocked:default").Return(false, nil)

	for i := 0; i < tknCfg.MaxRequests; i++ {
		suite.storage.On("Increment", ctx, "ratelimit:req_qnt:default", suite.ttl).Return(i+1, nil).Once()
		err := suite.limiter.Limit("127.0.0.1", tknCfg.Name)
		assert.NoError(suite.T(), err)
	}

	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:default", suite.ttl).Return(tknCfg.MaxRequests+1, nil).Once()
	suite.storage.On("Set", ctx, "ratelimit:blocked:default", tknCfg.Cooldown).Return(nil)
	err := suite.limiter.Limit("127.0.0.1", tknCfg.Name)
	assert.EqualError(suite.T(), err, limiter.ErrLimitedAccess.Error())

	// Esperar o tempo de bloqueio
	time.Sleep(time.Duration(tknCfg.Cooldown) * time.Second)

	suite.storage.On("IsBlocked", ctx, "ratelimit:blocked:default").Return(false, nil)
	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:default", suite.ttl).Return(1, nil).Once()
	err = suite.limiter.Limit("127.0.0.1", tknCfg.Name)
	assert.NoError(suite.T(), err)
}

func (suite *LimiterSuite) TestLimit_IP_Success() {
	ctx := context.Background()
	suite.storage.On("IsBlocked", ctx, "ratelimit:blocked:127.0.0.1").Return(false, nil)
	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:127.0.0.1", suite.ttl).Return(1, nil)
	err := suite.limiter.Limit("127.0.0.1", "")
	assert.NoError(suite.T(), err)
}

func (suite *LimiterSuite) TestLimit_ExceedIpMaxRequests_Error() {
	ctx := context.Background()
	suite.storage.On("IsBlocked", ctx, "ratelimit:blocked:127.0.0.1").Return(false, nil)

	for i := 0; i < suite.maxRequests; i++ {
		suite.storage.On("Increment", ctx, "ratelimit:req_qnt:127.0.0.1", suite.ttl).Return(i+1, nil).Once()
		err := suite.limiter.Limit("127.0.0.1", "")
		assert.NoError(suite.T(), err)
	}

	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:127.0.0.1", suite.ttl).Return(suite.maxRequests+1, nil).Once()
	suite.storage.On("Set", ctx, "ratelimit:blocked:127.0.0.1", suite.cooldown).Return(nil).Once()

	err := suite.limiter.Limit("127.0.0.1", "")
	assert.EqualError(suite.T(), err, limiter.ErrLimitedAccess.Error())
}

func (suite *LimiterSuite) TestLimit_IpBlockingDuration() {
	ctx := context.Background()
	suite.storage.On("IsBlocked", ctx, "ratelimit:blocked:127.0.0.1").Return(false, nil)

	for i := 0; i < suite.maxRequests; i++ {
		suite.storage.On("Increment", ctx, "ratelimit:req_qnt:127.0.0.1", suite.ttl).Return(i+1, nil).Once()
		err := suite.limiter.Limit("127.0.0.1", "")
		assert.NoError(suite.T(), err)
	}

	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:127.0.0.1", suite.ttl).Return(suite.maxRequests+1, nil).Once()
	suite.storage.On("Set", ctx, "ratelimit:blocked:127.0.0.1", suite.cooldown).Return(nil).Once()
	err := suite.limiter.Limit("127.0.0.1", "")
	assert.EqualError(suite.T(), err, limiter.ErrLimitedAccess.Error())

	// Esperar o tempo de bloqueio
	time.Sleep(time.Duration(suite.cooldown) * time.Second)

	suite.storage.On("Increment", ctx, "ratelimit:req_qnt:127.0.0.1", suite.ttl).Return(1, nil).Once()
	err = suite.limiter.Limit("127.0.0.1", "")
	assert.NoError(suite.T(), err)
}

func TestLimiterSuite(t *testing.T) {
	suite.Run(t, new(LimiterSuite))
}
