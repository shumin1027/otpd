package badger

import (
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewLogger(logger *zap.Logger) *Logger {
	l := logger.WithOptions(zap.Fields(zap.String("mod", "badger")))
	return &Logger{l.Sugar()}
}

// impl interface badger.Logger

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args)
}
func (l *Logger) Warningf(template string, args ...interface{}) {
	l.logger.Warnf(template, args)
}
func (l *Logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args)
}
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args)
}
