package logger

import (
	"log"
	"strings"
)

// Config Logger config.
type Config struct {
	// Filenames filename to storage log, could be stdout, stderr and files.
	Filenames []string `json:"filename" yaml:"filename"`
	// MaxSize log file size, if filename is a path to file, default to 100 (mb).
	MaxSize int `json:"max_size" yaml:"max_size"`
	// MaxAge log file retain time, if filename is a path to file, default to 0, not retain (day).
	MaxAge int `json:"max_age" yaml:"max_age"`
	// MaxBackups max files to storage, if filename is a path to file, , default to 0, not retain.
	MaxBackups int `json:"max_backups" yaml:"max_backups"`
	// LocalTime using local time, if filename is a path to file.
	LocalTime bool `json:"local_time" yaml:"local_time"`
	// Compress rotated file is compressed, if filename is a path to file.
	Compress bool `json:"compress" yaml:"compress"`
	// LogLevel support debug, info, warn, error, dpanic, panic, fatal, do not care uppgercase or lowwercase.
	LogLevel string `json:"log_level" yaml:"log_level"`
	// Encoder encoder log format to store or print, support json and console.
	Encoder string `json:"encoder" yaml:"encoder"`
}

// Build build config to fix all empty values.
func (c *Config) Build() {
	if len(c.Filenames) == 0 {
		c.Filenames = []string{LOGGER_FILE_STDERR}
	}
	if c.MaxSize < 0 {
		c.MaxSize = 0
	}
	if c.MaxAge < 0 {
		c.MaxAge = 0
	}
	if c.MaxBackups < 0 {
		c.MaxBackups = 0
	}
	switch strings.ToLower(c.LogLevel) {
	case "debug", "info", "warn", "error", "dpanic", "panic", "fatal", "":
		c.LogLevel = strings.ToLower(c.LogLevel)
	default:
		log.Printf("unknown log level: %s, set to info as default", c.LogLevel)
		c.LogLevel = ""
	}
	switch strings.ToLower(c.Encoder) {
	case LOGGER_ENCODER_CONSOLE, LOGGER_ENCODER_JSON:
		c.Encoder = strings.ToLower(c.Encoder)
	case "":
		c.Encoder = LOGGER_ENCODER_CONSOLE
	default:
		log.Printf("unknown encoder: %s, set to console as default", c.Encoder)
		c.Encoder = LOGGER_ENCODER_CONSOLE
	}
}
