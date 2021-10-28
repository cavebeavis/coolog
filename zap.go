package coolog

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
    https://pkg.go.dev/go.uber.org/zap#Config
    https://play.golang.org/p/8PTI83OA8XZ
    // For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
*/

/*
	https://pkg.go.dev/go.uber.org/zap@v1.17.0/zapcore#EncoderConfig
	type EncoderConfig struct {
		// Set the keys used for each log entry. If any key is empty, that portion
		// of the entry is omitted.
		MessageKey    string `json:"messageKey" yaml:"messageKey"`
		LevelKey      string `json:"levelKey" yaml:"levelKey"`
		TimeKey       string `json:"timeKey" yaml:"timeKey"`
		NameKey       string `json:"nameKey" yaml:"nameKey"`
		CallerKey     string `json:"callerKey" yaml:"callerKey"`
		FunctionKey   string `json:"functionKey" yaml:"functionKey"`
		StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
		LineEnding    string `json:"lineEnding" yaml:"lineEnding"`
		// Configure the primitive representations of common complex types. For
		// example, some users may want all time.Times serialized as floating-point
		// seconds since epoch, while others may prefer ISO8601 strings.
		EncodeLevel    LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
		EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
		EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
		EncodeCaller   CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
		// Unlike the other primitive type encoders, EncodeName is optional. The
		// zero value falls back to FullNameEncoder.
		EncodeName NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
		// Configures the field separator used by the console encoder. Defaults
		// to tab.
		ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"`
	}
*/

// ZapLogger is the struct which implements the Logger interface which allows access to the
// "go.uber.org/zap" style of logging.
type ZapLogger struct {
	Logger *zap.Logger
	Level  string
	Atom   zap.AtomicLevel
}

// NewZapLogger provides a newly initialized Zap Logger which implements Logger interface.
//
// level is "trace", "debug", "info", "warn", "error" which will be the minimum level
// this Logger will log in the file or console. If level is empty, it will default to
// "info" level.
//
// logLocations s where you want the log file stored. If this is empty, it will default
// to stdout.
//
// logType is "json" or "text". If this is empty, it will default to "json".
func NewZapLogger(level string, logLocations []string, logType string) (*ZapLogger, error) {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "@timestamp", // For ElasticSearch: https://www.elastic.co/guide/en/ecs/current/ecs-base.html#field-timestamp
		//NameKey       string `json:"nameKey" yaml:"nameKey"`
		CallerKey:     "caller",
		FunctionKey:   "func",
		StacktraceKey: "stack",
		LineEnding:    ",\n",
		// Configure the primitive representations of common complex types. For
		// example, some users may want all time.Times serialized as floating-point
		// seconds since epoch, while others may prefer ISO8601 strings.
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339Nano), // TODO: make this user configurable.
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		// Unlike the other primitive type encoders, EncodeName is optional. The
		// zero value falls back to FullNameEncoder.
		//EncodeName NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
		// Configures the field separator used by the console encoder. Defaults
		// to tab.
		//ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"`
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

	for i, location := range logLocations {
		if location == "console" || location == "" {
			logLocations[i] = "stdout"
		}
	}

	switch logType {
	case "txt", "text", "plain":
		logType = "console"
	case "console":
		// fall through
	default:
		logType = "json"
	}

	cfg := zap.Config{
		EncoderConfig:    encoderCfg,
		Level:            atom,
		Encoding:         logType,
		OutputPaths:      logLocations,
		ErrorOutputPaths: logLocations,
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.WithOptions(zap.AddCallerSkip(1))

	return &ZapLogger{Logger: logger, Atom: atom, Level: level}, nil

}

// Print is the main method, and implements interface Logger.
func (z *ZapLogger) Print(level string, msg string, data ...map[string]interface{}) error {
	//TODO: get rid of Any and do this with true zap fields from the calling function...
	var fields []zap.Field

	if len(data) > 0 {
		for _, dataMap := range data {
			if dataMap == nil {
				continue
			}

			for k, v := range dataMap {
				fields = append(fields, zap.Any(k, v))
			}
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
	case "fatal":
		z.Logger.Fatal(msg, fields...)
	case "panic":
		z.Logger.Panic(msg, fields...)
	default:
		z.Logger.Error(msg, fields...)
	}

	return nil
}

// Close will perform the Logger.Sync() specifically for zap.
func (z *ZapLogger) Close() error {
	return z.Logger.Sync()
}
