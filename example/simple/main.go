package main

import (
	"fmt"

	"github.com/cavebeavis/coolog"
)

func main() {
	fmt.Println("Starting simple logging example...")

	logLocations := []string{"stdout"}
	logLevel := "info"

	// this is the zap logger flavor, but the point is any logger can be used
	// without having to rewire iNeedsLogger below as long as the new logger
	// implements the coolog.Logger interface...
	zapLogger, err := coolog.NewZapLogger(logLevel, logLocations, "text")
	if err != nil {
		fmt.Println("well this is seriously embarrassing, zap:", err)
		return
	}
	defer zapLogger.Close()

	iNeedsLogger(zapLogger)

	logrusLogger, err := coolog.NewLogrusLogger(logLevel, logLocations[0], "text")
	if err != nil {
		fmt.Println("well this is seriously embarrassing, logrus:", err)
		return
	}

	iNeedsLogger(logrusLogger)

	fmt.Println("Finished")
}

func iNeedsLogger(log coolog.Logger) {
	// we do lots of stuff and need to log
	log.Print("info", "I am info, hear me roar!")

	// maybe we have data and fields...
	log.Print("error", "I am an error with context data", map[string]interface{}{
		"field1": "you really did it now",
		"field2": 666,
		"field3": []string{"slice", "is", "ok", "too"},
	})

	// or trace logs but we don't want to see these unless the Logger is in trace mode
	log.Print("trace", "You will only see me if trace or debug (for zap implementation) is enabled")
}