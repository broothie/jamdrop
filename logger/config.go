package logger

import (
	"io"
	"os"
	"time"
)

type config struct {
	Level      Level
	Writer     io.Writer
	BufferSize uint
	TimeFormat string
}

func defaultConfig() config {
	return config{
		Level:      Debug,
		Writer:     os.Stdout,
		BufferSize: 64,
		TimeFormat: time.RFC3339Nano,
	}
}

type Configurer func(*config)

func ConfigureLevel(level Level) Configurer {
	return func(c *config) { c.Level = level }
}

func ConfigureWriter(writer io.Writer) Configurer {
	return func(c *config) { c.Writer = writer }
}

func ConfigureBufferSize(size uint) Configurer {
	return func(c *config) { c.BufferSize = size }
}

func ConfigureTimeFormat(timeFormat string) Configurer {
	return func(c *config) { c.TimeFormat = timeFormat }
}
