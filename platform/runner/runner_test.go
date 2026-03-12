package runner

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/ovya/ogl/oglcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConfig
type MockConfig struct {
	mock.Mock
}

type stringer string

func (s stringer) String() string { return string(s) }

func (m *MockConfig) GetAppEnv() fmt.Stringer {
	args := m.Called()
	return args.Get(0).(fmt.Stringer)
}
func (m *MockConfig) GetAppName() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockConfig) GetServerPort() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockConfig) GetServerHost() string {
	args := m.Called()
	return args.String(0)
}
func (m *MockConfig) GetDatabaseURL() string {
	args := m.Called()
	return args.String(0)
}

// MockPinger
type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockModule
type MockModule struct {
	mock.Mock
}

func (m *MockModule) RegisterRoutes(mux *http.ServeMux) {
	m.Called(mux)
}
func (m *MockModule) StartWorkers(ctx context.Context) error {
	args := m.Called(ctx)
	// Block until context is done to simulate worker
	<-ctx.Done()
	return args.Error(0)
}
func (m *MockModule) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNew(t *testing.T) {
	cfg := new(MockConfig)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	db := new(MockPinger)

	app := New(cfg, logger, db, nil)
	assert.NotNil(t, app)
	assert.Equal(t, cfg, app.config)
	assert.Equal(t, logger, app.logger)
	assert.Equal(t, db, app.db)
}

func TestApp_Run(t *testing.T) {
	// Setup
	cfg := new(MockConfig)
	cfg.On("GetAppName").Return("test-app")
	cfg.On("GetServerPort").Return(":0") // Random port

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	db := new(MockPinger)

	mod := new(MockModule)
	mod.On("RegisterRoutes", mock.Anything).Return()
	mod.On("StartWorkers", mock.Anything).Return(nil)
	mod.On("Close").Return(nil)

	app := New(cfg, logger, db, []oglcore.Module{mod})

	// Run
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	go func() {
		errChan <- app.Run(ctx)
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context to stop
	cancel()

	// Wait for return
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("App.Run did not return after context cancellation")
	}

	mod.AssertExpectations(t)
	cfg.AssertExpectations(t)
}
