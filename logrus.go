package coolog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// LogrusLogger is our struct to implement the Logger interface.
type LogrusLogger struct {
	Logger      *logrus.Logger
	LogrusLevel logrus.Level
	Level       string
}

// NewLogrusLogger is a convenience function to create a logrus Logger which implements
// Logger interface.
//
// logLevel is "trace", "debug", "info", "warn", "error", "fatal", or "panic" and will be
// the minimum level the Logger will log in the file or console. If level is empty, it
// will default to "info" level.
//
// logFilename is where you want the log file stored. If this is empty, it will default
// to stdout.
//
// logType is "json" or "text". If this is empty, it will default to "json".
func NewLogrusLogger(logLevel, logFilename, logType string) (*LogrusLogger, error) {
	logger := logrus.New()

	logger.SetReportCaller(true) // so we can get the filename, function, and line number.

	callerPrettyfier := func(f *runtime.Frame) (string, string) {
		//filename := path.Base(f.File)
		ptr, filename, ln, _ := runtime.Caller(7)
		frames := runtime.CallersFrames([]uintptr{ptr})
		frame, _ := frames.Next()

		funcParts := strings.Split(frame.Function, ".")

		return fmt.Sprintf("%s()", funcParts[len(funcParts)-1]), fmt.Sprintf("%s:%d", filename, ln)
	}

	// TODO: add optional opts ...map[string]string to the function to allow user modification.
	fieldMap := logrus.FieldMap{
		logrus.FieldKeyMsg:         "msg",
		logrus.FieldKeyLevel:       "lvl",
		logrus.FieldKeyTime:        "@timestamp", // For ElasticSearch: https://www.elastic.co/guide/en/ecs/master/ecs-base.html#field-timestamp
		logrus.FieldKeyLogrusError: "logrusError",
		logrus.FieldKeyFunc:        "func",
		logrus.FieldKeyFile:        "filepath",
	}

	switch logType {
	case "txt", "text", "plain":
		logger.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: callerPrettyfier,
			TimestampFormat:  time.RFC3339Nano, // TODO: another user configurable field?
			FieldMap:         fieldMap,
		}
	default:
		logger.Formatter = &logrus.JSONFormatter{
			CallerPrettyfier: callerPrettyfier,
			TimestampFormat:  time.RFC3339Nano, // TODO: another user configurable field?
			FieldMap:         fieldMap,
		}
	}

	logger.SetOutput(os.Stdout)
	if logFilename != "" && logFilename != "console" {
		logger.Println("logging will now be saved in ", logFilename)

		f, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return nil, err
		}

		log.SetOutput(f)
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Error(err)
		level = logrus.InfoLevel
		logger.Info("setting log to info level")
	}

	logger.SetLevel(level)

	return &LogrusLogger{
		Logger: logger,
		Level:  logger.Level.String(),
	}, nil
}

// Print is the main method, and implements interface Logger.
func (l *LogrusLogger) Print(level string, msg string, data ...map[string]interface{}) error {
	fields := make(logrus.Fields)

	if len(data) > 0 {
		for _, dataMap := range data {
			if dataMap == nil {
				continue
			}

			for k, v := range dataMap {
				fields[k] = v
			}
		}
	}

	switch level {
	case "trace":
		l.Logger.WithFields(fields).Trace(msg)
	case "debug":
		l.Logger.WithFields(fields).Debug(msg)
	case "info":
		l.Logger.WithFields(fields).Info(msg)
	case "warn":
		l.Logger.WithFields(fields).Warn(msg)
	case "error":
		l.Logger.WithFields(fields).Error(msg)
	case "fatal":
		l.Logger.WithFields(fields).Fatal(msg)
	case "panic":
		l.Logger.WithFields(fields).Panic(msg)
	default:
		l.Logger.WithFields(fields).Error(msg)
	}

	return nil
}