package limiter

import (
	"testing"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/configs"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/mocks"
	"github.com/stretchr/testify/assert"
)

type TokenConfig = configs.TokenConfig

// type LimiterTestSuite struct {
// 	suite.Suite
// 	storageMock *mocks.StorageMock
// }

func TestLimit_IpAccessAllowed(t *testing.T) {
	storage := &mocks.StorageMock{}
	cooldown := 3
	ttl := 1
	maxRequests := 10
	tokens := make(map[string]configs.TokenConfig)

	ip := "127.0.0.1"
	token := ""
	term := ip

	storage.On("IsBlocked", term).Return(false, nil)
	storage.On("IncrementCounter", term, ttl).Return(1, nil)

	limiter := NewLimiter(storage, cooldown, maxRequests, ttl, tokens)

	err := limiter.Limit(ip, token)
	assert.Nil(t, err)

	storage.AssertExpectations(t)
}

func TestLimit_IpAccessNotAllowed(t *testing.T) {
	storage := &mocks.StorageMock{}
	cooldown := 3
	ttl := 1
	maxRequests := 10
	tokens := make(map[string]configs.TokenConfig)

	ip := "127.0.0.1"
	token := ""
	term := ip

	storage.On("IsBlocked", term).Return(false, nil)
	storage.On("IncrementCounter", term, ttl).Return(11, nil)
	storage.On("RegisterBlock", term, cooldown).Once().Return(nil)

	limiter := NewLimiter(storage, cooldown, maxRequests, ttl, tokens)

	err := limiter.Limit(ip, token)
	assert.NotNil(t, err)
	assert.Equal(t, ErrLimitedAccess, err, "access blocked")

	storage.AssertExpectations(t)
}

func TestLimit_TokenAccessAllowed(t *testing.T) {
	storage := &mocks.StorageMock{}
	cooldown := 3
	ttl := 1

	tokens := make(map[string]configs.TokenConfig)
	tokens["tkn_123"] = TokenConfig{Name: "tkn_123", MaxRequests: 20, Cooldown: 2}
	tokens["tkn_456"] = TokenConfig{Name: "tkn_456", MaxRequests: 30, Cooldown: 1}

	ip := "127.0.0.1"
	token := "tkn_123"
	term := token
	maxRequests := tokens[term].MaxRequests

	storage.On("IsBlocked", term).Return(false, nil)
	storage.On("IncrementCounter", term, ttl).Return(1, nil)

	limiter := NewLimiter(storage, cooldown, maxRequests, ttl, tokens)

	err := limiter.Limit(ip, token)
	assert.Nil(t, err)

	storage.AssertExpectations(t)
}

func TestLimit_TokenAccessNotAllowed(t *testing.T) {
	storage := &mocks.StorageMock{}
	cooldown := 3
	ttl := 1

	tokens := make(map[string]configs.TokenConfig)
	tokens["tkn_123"] = TokenConfig{Name: "tkn_123", MaxRequests: 20, Cooldown: 2}
	tokens["tkn_456"] = TokenConfig{Name: "tkn_456", MaxRequests: 30, Cooldown: 1}

	ip := "127.0.0.1"
	token := "tkn_123"
	term := token
	maxRequests := tokens[term].MaxRequests

	storage.On("IsBlocked", term).Return(false, nil)
	storage.On("IncrementCounter", term, ttl).Return(maxRequests+1, nil)
	storage.On("RegisterBlock", term, tokens[term].Cooldown).Once().Return(nil)

	limiter := NewLimiter(storage, cooldown, maxRequests, ttl, tokens)

	err := limiter.Limit(ip, token)
	assert.NotNil(t, err)
	assert.Equal(t, ErrLimitedAccess, err, "access blocked")

	storage.AssertExpectations(t)
}

func TestLimit_AccessNotAllowed_Already_Blocked(t *testing.T) {
	storage := &mocks.StorageMock{}
	cooldown := 3
	ttl := 1

	tokens := make(map[string]configs.TokenConfig)
	tokens["tkn_123"] = TokenConfig{Name: "tkn_123", MaxRequests: 20, Cooldown: 2}
	tokens["tkn_456"] = TokenConfig{Name: "tkn_456", MaxRequests: 30, Cooldown: 1}

	ip := "127.0.0.1"
	token := "tkn_123"
	term := token

	maxRequests := tokens[term].MaxRequests

	storage.On("IsBlocked", term).Return(true, nil)
	limiter := NewLimiter(storage, cooldown, maxRequests, ttl, tokens)

	err := limiter.Limit(ip, token)
	assert.NotNil(t, err)
	assert.Equal(t, ErrLimitedAccess, err, "access blocked")

	storage.AssertExpectations(t)
}

func TestLimit_TokenAccess(t *testing.T) {
	storage := &mocks.StorageMock{}
	cooldown := 3
	ttl := 1

	tokens := make(map[string]TokenConfig)
	tokens["tkn_123"] = TokenConfig{Name: "tkn_123", MaxRequests: 2, Cooldown: 2}
	tokens["tkn_456"] = TokenConfig{Name: "tkn_456", MaxRequests: 3, Cooldown: 1}

	ip := "127.0.0.1"
	token := "tkn_123"
	term := token

	maxRequests := tokens[term].MaxRequests

	limiter := NewLimiter(storage, cooldown, maxRequests, ttl, tokens)
	storage.On("IsBlocked", term).Return(false, nil).Times(maxRequests)

	for i := 1; i <= maxRequests; i++ {
		storage.On("IncrementCounter", term, ttl).Return(i, nil).Once()
		err := limiter.Limit(ip, token)
		assert.Nil(t, err)
	}

	storage.On("IsBlocked", term).Return(false, nil).Once()
	storage.On("IncrementCounter", term, ttl).Once().Return(maxRequests+1, nil)
	storage.On("RegisterBlock", term, tokens[term].Cooldown).Return(nil).Once()
	storage.On("IsBlocked", term).Return(true, nil)

	// Simulando requests após o limite permitido
	for i := 0; i < 50; i++ {
		err := limiter.Limit(term, token)
		assert.NotNil(t, err)
		assert.Equal(t, ErrLimitedAccess, err, "acesso limitado")
	}

	storage.AssertExpectations(t)
}

func TestLimit_IpAccess(t *testing.T) {
	storage := &mocks.StorageMock{}
	cooldown := 3
	ttl := 1

	tokens := make(map[string]TokenConfig)

	ip := "127.0.0.1"
	token := "tkn_123"
	term := ip

	maxRequests := 10

	limiter := NewLimiter(storage, cooldown, maxRequests, ttl, tokens)
	storage.On("IsBlocked", term).Return(false, nil).Times(maxRequests)

	for i := 1; i <= maxRequests; i++ {
		storage.On("IncrementCounter", term, ttl).Return(i, nil).Once()
		err := limiter.Limit(ip, token)
		assert.Nil(t, err)
	}

	storage.On("IsBlocked", term).Return(false, nil).Once()
	storage.On("IncrementCounter", term, ttl).Once().Return(maxRequests+1, nil)
	storage.On("RegisterBlock", term, cooldown).Return(nil).Once()
	storage.On("IsBlocked", term).Return(true, nil)

	// Simulando requests após o limite permitido
	for i := 0; i < 50; i++ {
		err := limiter.Limit(term, token)
		assert.NotNil(t, err)
		assert.Equal(t, ErrLimitedAccess, err, "acesso limitado")
	}

	storage.AssertExpectations(t)
}

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// type LimiterTestSuite struct {
// 	suite.Suite
// 	controller         *gomock.Controller
// 	context            context.Context
// 	storageAdapterMock *mocks.MockRateLimitStorageAdapter
// }

// func TestRateLimiterTestSuite(t *testing.T) {
// 	suite.Run(t, new(RateLimiterTestSuite))
// }

// func (s *RateLimiterTestSuite) SetupTest() {
// 	s.controller = gomock.NewController(s.T())
// 	s.context = context.Background()
// 	s.storageAdapterMock = mocks.NewMockRateLimitStorageAdapter(s.controller)
// }

// func (s *RateLimiterTestSuite) TestCheckRateLimit_AccessAllowed() {
// 	context := s.context
// 	keyType := "IP"
// 	key := "127.0.0.1"
// 	config := &RateLimiterConfig{
// 		IP: &RateLimiterRateConfig{
// 			MaxRequestsPerSecond:  10,
// 			BlockTimeMilliseconds: 100,
// 		},
// 	}

// 	s.storageAdapterMock.EXPECT().
// 		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

// 	s.storageAdapterMock.EXPECT().
// 		IncrementAccesses(context, keyType, key, gomock.Any()).Return(true, int64(1), nil).Times(1)

// 	config.StorageAdapter = s.storageAdapterMock

// 	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
// 	assert.Nil(s.T(), err)
// 	assert.Nil(s.T(), returnedBlock)
// }

// func (s *RateLimiterTestSuite) TestCheckRateLimit_AccessDenied() {
// 	context := s.context
// 	keyType := "IP"
// 	key := "127.0.0.1"
// 	config := &RateLimiterConfig{
// 		IP: &RateLimiterRateConfig{
// 			MaxRequestsPerSecond:  10,
// 			BlockTimeMilliseconds: 100,
// 		},
// 	}
// 	block := time.Now().Add(time.Millisecond * 100)

// 	s.storageAdapterMock.EXPECT().
// 		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

// 	s.storageAdapterMock.EXPECT().
// 		IncrementAccesses(context, keyType, key, gomock.Any()).Return(false, int64(10), nil).Times(1)

// 	s.storageAdapterMock.EXPECT().
// 		AddBlock(context, keyType, key, config.IP.BlockTimeMilliseconds).Return(&block, nil).Times(1)

// 	config.StorageAdapter = s.storageAdapterMock

// 	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
// 	assert.Nil(s.T(), err)
// 	assert.Equal(s.T(), block, *returnedBlock)
// }

// func (s *RateLimiterTestSuite) TestCheckRateLimit_AlreadyBlocked() {
// 	context := s.context
// 	keyType := "IP"
// 	key := "127.0.0.1"
// 	config := &RateLimiterConfig{
// 		IP: &RateLimiterRateConfig{
// 			MaxRequestsPerSecond:  10,
// 			BlockTimeMilliseconds: 100,
// 		},
// 	}
// 	block := time.Now().Add(time.Millisecond * 100)

// 	s.storageAdapterMock.EXPECT().
// 		GetBlock(context, keyType, key).Return(&block, nil).Times(1)

// 	config.StorageAdapter = s.storageAdapterMock

// 	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
// 	assert.Nil(s.T(), err)
// 	assert.Equal(s.T(), block, *returnedBlock)
// }

// func (s *RateLimiterTestSuite) TestCheckRateLimit_EmptyKey() {
// 	context := s.context
// 	keyType := "IP"
// 	key := ""
// 	config := &RateLimiterConfig{
// 		IP: &RateLimiterRateConfig{
// 			MaxRequestsPerSecond:  10,
// 			BlockTimeMilliseconds: 100,
// 		},
// 	}

// 	config.StorageAdapter = s.storageAdapterMock

// 	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
// 	assert.Nil(s.T(), err)
// 	assert.Nil(s.T(), returnedBlock)
// }

// func (s *RateLimiterTestSuite) TestCheckRateLimit_GetBlockError() {
// 	context := s.context
// 	keyType := "IP"
// 	key := "127.0.0.1"
// 	config := &RateLimiterConfig{
// 		IP: &RateLimiterRateConfig{
// 			MaxRequestsPerSecond:  10,
// 			BlockTimeMilliseconds: 100,
// 		},
// 	}

// 	s.storageAdapterMock.EXPECT().
// 		GetBlock(context, keyType, key).Return(nil, errors.New("error")).Times(1)

// 	config.StorageAdapter = s.storageAdapterMock

// 	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
// 	assert.NotNil(s.T(), err)
// 	assert.Nil(s.T(), returnedBlock)
// }

// func (s *RateLimiterTestSuite) TestCheckRateLimit_IncrementAccessesError() {
// 	context := s.context
// 	keyType := "IP"
// 	key := "127.0.0.1"
// 	config := &RateLimiterConfig{
// 		IP: &RateLimiterRateConfig{
// 			MaxRequestsPerSecond:  10,
// 			BlockTimeMilliseconds: 100,
// 		},
// 	}

// 	s.storageAdapterMock.EXPECT().
// 		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

// 	s.storageAdapterMock.EXPECT().
// 		IncrementAccesses(context, keyType, key, gomock.Any()).Return(false, int64(1), errors.New("error")).Times(1)

// 	config.StorageAdapter = s.storageAdapterMock

// 	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
// 	assert.NotNil(s.T(), err)
// 	assert.Nil(s.T(), returnedBlock)
// }

// func (s *RateLimiterTestSuite) TestCheckRateLimit_AddBlockError() {
// 	context := s.context
// 	keyType := "IP"
// 	key := "127.0.0.1"
// 	config := &RateLimiterConfig{
// 		IP: &RateLimiterRateConfig{
// 			MaxRequestsPerSecond:  10,
// 			BlockTimeMilliseconds: 100,
// 		},
// 	}

// 	s.storageAdapterMock.EXPECT().
// 		GetBlock(context, keyType, key).Return(nil, nil).Times(1)

// 	s.storageAdapterMock.EXPECT().
// 		IncrementAccesses(context, keyType, key, gomock.Any()).Return(false, int64(10), nil).Times(1)

// 	s.storageAdapterMock.EXPECT().
// 		AddBlock(context, keyType, key, config.IP.BlockTimeMilliseconds).Return(nil, errors.New("error")).Times(1)

// 	config.StorageAdapter = s.storageAdapterMock

// 	returnedBlock, err := checkRateLimit(context, keyType, key, config, config.IP)
// 	assert.NotNil(s.T(), err)
// 	assert.Nil(s.T(), returnedBlock)
// }
