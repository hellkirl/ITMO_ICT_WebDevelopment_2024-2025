// dbcontext/dbcontext.go
package dbcontext

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"lr2/pkg/logger"
	"time"
)

type DbContext struct {
	pool         *pgxpool.Pool
	logger       *logger.ZapLogger
	connURL      string
	watchCtx     context.Context
	cancelWatch  context.CancelFunc
	watchStopped chan struct{}
}

func NewDbContext(connURL string, logger *logger.ZapLogger) (*DbContext, error) {
	config, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		logger.Error("Failed to parse connection URL", zap.Error(err))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("Couldn't set up pool connection with database", zap.Error(err))
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		logger.Error("Couldn't ping the database", zap.Error(err))
		pool.Close()
		return nil, err
	}

	watchCtx, cancel := context.WithCancel(context.Background())
	dc := &DbContext{
		pool:         pool,
		logger:       logger,
		connURL:      connURL,
		watchCtx:     watchCtx,
		cancelWatch:  cancel,
		watchStopped: make(chan struct{}),
	}

	go dc.watchConnection()

	return dc, nil
}

func (dc *DbContext) Pool() (*pgxpool.Pool, error) {
	if dc.pool == nil {
		return nil, fmt.Errorf("failed to init db pool")
	}
	return dc.pool, nil
}

func (dc *DbContext) watchConnection() {
	defer close(dc.watchStopped)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-dc.watchCtx.Done():
			dc.logger.Info("Connection watcher stopped")
			return
		case <-ticker.C:
			dc.logger.Debug("Pinging database to check connection health")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := dc.pool.Ping(ctx)
			cancel()
			if err != nil {
				dc.logger.Error("Database connection lost, attempting to reconnect", zap.Error(err))
				if reconnectErr := dc.reconnect(); reconnectErr != nil {
					dc.logger.Error("Reconnection attempt failed", zap.Error(reconnectErr))
				} else {
					dc.logger.Info("Reconnected to the database successfully")
				}
			} else {
				dc.logger.Debug("Database connection is healthy")
			}
		}
	}
}

func (dc *DbContext) reconnect() error {
	dc.pool.Close()

	newPool, err := pgxpool.New(context.Background(), dc.connURL)
	if err != nil {
		dc.logger.Error("Failed to create new connection pool during reconnection", zap.Error(err))
		return err
	}

	if err := newPool.Ping(context.Background()); err != nil {
		dc.logger.Error("Failed to ping the database with the new connection pool", zap.Error(err))
		newPool.Close()
		return err
	}

	dc.pool = newPool
	return nil
}

func (dc *DbContext) Close() {
	dc.cancelWatch()
	<-dc.watchStopped
	dc.pool.Close()
	dc.logger.Info("DbContext closed")
}
