package main

import (
	"SocialMediaApp/internal/auth"
	"SocialMediaApp/internal/ratelimiter"
	"SocialMediaApp/internal/store"
	"SocialMediaApp/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	//logger := zap.NewNop().Sugar()
	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStorage()

	testAuth := &auth.TestAuthenticator{}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
		rateLimiter:   rateLimiter,
	}
}

func executeRequest(req *http.Request, mux *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr

}
