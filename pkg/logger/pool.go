package logger

import (
	"errors"
	"sync"

	"go.uber.org/zap"
)

var (
	pool *LoggerPool
	// ErrLoggerExists logger is already exists.
	ErrLoggerExists = errors.New("logger already exists")
	// ErrLoggerExists logger is not exists.
	ErrLoggerNotExists = errors.New("logger does not exists")
)

// LoggerPool is a pool to manager logger for different business.
type LoggerPool struct {
	rwlock *sync.RWMutex
	pmap   map[string]*zap.Logger
}

// P global pool.
func P() *LoggerPool {
	return pool
}

func NewPool() *LoggerPool {
	return &LoggerPool{
		rwlock: &sync.RWMutex{},
		pmap:   make(map[string]*zap.Logger),
	}
}

// Add add a named logger to pool.
func (p *LoggerPool) Add(name string, config Config) (*zap.Logger, error) {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	if _, ok := p.pmap[name]; ok {
		return nil, ErrLoggerExists
	}
	p.pmap[name] = NewLogger(config)
	return p.pmap[name], nil
}

// Remove remove a named logger from pool.
func (p *LoggerPool) Remove(name string) error {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	_, ok := p.pmap[name]
	if !ok {
		return nil
	}
	delete(p.pmap, name)
	return nil
}

// Get get a named logger from pool.
func (p *LoggerPool) Get(name string) (*zap.Logger, error) {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()
	logger, ok := p.pmap[name]
	if !ok {
		return nil, ErrLoggerNotExists
	}
	return logger, nil
}
