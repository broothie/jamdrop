package logger

type Logger struct {
	config
	itemChan chan Item
}

func New(configurers ...Configurer) *Logger {
	cfg := defaultConfig()
	for _, configurer := range configurers {
		configurer(&cfg)
	}

	logger := &Logger{
		config:   cfg,
		itemChan: make(chan Item, cfg.BufferSize),
	}

	go logger.worker()
	return logger
}
