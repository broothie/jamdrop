package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type Fieldser interface {
	Fields() map[string]interface{}
}

type Fields map[string]interface{}

func (f Fields) Fields() map[string]interface{} {
	return f
}

func Field(key string, value interface{}) Fields {
	return Fields{key: value}
}

type Item struct {
	Level   Level
	Message string
	Fields  []Fieldser
	Time    time.Time
}

type Logger struct {
	Level    Level
	itemChan chan Item
}

func New(level Level) *Logger {
	logger := &Logger{Level: level, itemChan: make(chan Item)}
	go logger.worker()
	return logger
}

func (l *Logger) Log(level Level, message string, fields ...Fieldser) {
	now := time.Now().UTC()
	go func() { l.itemChan <- Item{Level: level, Message: message, Fields: fields, Time: now} }()
}

func (l *Logger) Err(err error, message string, fields ...Fieldser) {
	l.Error(message, append(fields, Fields{"error": err.Error()})...)
}

func (l *Logger) worker() {
	for item := range l.itemChan {
		if item.Level < l.Level {
			continue
		}

		payload := make(map[string]interface{})
		for _, fields := range item.Fields {
			for key, value := range fields.Fields() {
				if key != "" {
					payload[key] = value
				}
			}
		}

		payload["message"] = item.Message
		payload["level"] = item.Level.String()
		payload["time"] = item.Time
		encoded, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("failed to encode payload: payload: %v, error: %s\n", payload, err)
			return
		}

		fmt.Println(string(encoded))
	}
}
