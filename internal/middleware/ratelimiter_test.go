package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/limiter"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/middleware"
	"github.com/rhbarauna/goexpert-desafio-rate-limiter/internal/mocks"
)

type RateLimiterSuite struct {
	suite.Suite
	middleware.RateLimiter
	limiterMock *mocks.LimiterMock
}

func (suite *RateLimiterSuite) SetupTest() {
	suite.limiterMock = new(mocks.LimiterMock)
	suite.RateLimiter = middleware.NewRateLimiter(suite.limiterMock)
}

func (suite *RateLimiterSuite) TestValidRequest() {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(suite.T(), err)

	suite.limiterMock.On("Limit", mock.Anything, mock.Anything).Once().Return(nil)

	rr := httptest.NewRecorder()
	suite.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rr, req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.limiterMock.AssertExpectations(suite.T())
}

func (suite *RateLimiterSuite) TestInvalidRequest() {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(suite.T(), err)

	suite.limiterMock.On("Limit", mock.Anything, mock.Anything).Return(limiter.ErrLimitedAccess)

	rr := httptest.NewRecorder()

	suite.RateLimiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rr, req)

	assert.Error(suite.T(), limiter.ErrLimitedAccess, rr.Result().Body)
	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)

	expectedMessage := "You have reached the maximum number of requests or actions allowed within a certain time frame."
	assert.Contains(suite.T(), rr.Body.String(), expectedMessage)
	suite.limiterMock.AssertExpectations(suite.T())
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterSuite))
}
