package coolog

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LogrusLogger is our struct to implement the Logger interface.
type LogrusLogger struct {
    Logger      *logrus.Logger
	LogrusLevel logrus.Level
    Level       string
}

// NewLogrusLogger is a convenience function to create the LogrusLogger.
func NewLogrusLogger(level string, logLocation string) (*LogrusLogger, error) {
    logger := logrus.New()

	switch logLocation {
	case "console", "stdout":
		// fall through
	default:
		f, err := os.Create(logLocation)
		if err != nil {
			return nil, err
		}

		logger.SetOutput(f)
	}
	
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrusLevel = logrus.InfoLevel
	}

    logger.SetLevel(logrusLevel)

    return &LogrusLogger{Logger: logger, LogrusLevel: logrusLevel, Level: level}, nil
}

// Print is the main method, and implements interface Logger.
func (l *LogrusLogger) Print(level string, msg string, data ...map[string]interface{}) error {
    fields := make(logrus.Fields)

	for _, d := range data {
		for k, v := range d {
			fields[k] = v
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
