package coolog

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// ZapLogger is our struct to implement the Logger interface.
type ZapLogger struct {
    Logger  *zap.Logger
    Level   string
    Atom    zap.AtomicLevel
}

// NewZapLogger is a convenience function to create the ZapLogger.
func NewZapLogger(level string, logLocations []string) (*ZapLogger, error) {
    encoderCfg := zapcore.EncoderConfig{
        MessageKey:    "msg",
        LevelKey:      "level",
        TimeKey:       "@ts",
        CallerKey:     "caller",
        FunctionKey:   "func",
        StacktraceKey: "stack",
        LineEnding:    ",\n",  
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000000 MST"),
        EncodeDuration: zapcore.NanosDurationEncoder,
        EncodeCaller:   zapcore.FullCallerEncoder,
    }
    zapLevel := zap.InfoLevel
    switch level {
    case "trace", "debug":
        zapLevel = zap.DebugLevel
    case "info":
        zapLevel = zap.InfoLevel
    case "warn":
        zapLevel = zap.WarnLevel
    case "error":
        zapLevel = zap.ErrorLevel
    }
    atom := zap.NewAtomicLevelAt(zapLevel)
    cfg := zap.Config{
        EncoderConfig: encoderCfg,
        Level: atom,
        Encoding: "json",
        OutputPaths: logLocations,
        ErrorOutputPaths: logLocations,
    }
    logger, err := cfg.Build()
    if err != nil {
        return nil, err
    }
    return &ZapLogger{Logger: logger, Atom: atom, Level: level}, nil
}

// Print is the main method, and implements interface Logger.
func (z *ZapLogger) Print(level string, msg string, data ...map[string]interface{}) error {
    var fields []zap.Field
	for _, d := range data {
		for k, v := range d {
			field := zap.Any(k, v)
			fields = append(fields, field)
		}
	}
    
    switch level {
    case "trace", "debug":
        z.Logger.Debug(msg, fields...)
    case "info":
        z.Logger.Info(msg, fields...)
    case "warn":
        z.Logger.Warn(msg, fields...)
    case "error":
       z.Logger.Error(msg, fields...)
    default:
       z.Logger.Error(msg, fields...)
    }
    return nil
}

// Close is following the examples in zap -- may be a better way to
// do this since the caller of NewZapLogger needs to defer z.Close().
func (z *ZapLogger) Close() error {
    return z.Logger.Sync()
}