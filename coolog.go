package coolog

// Logger is influenced by gokit's log (https://pkg.go.dev/github.com/go-kit/kit@v0.10.0/log), but I
// did not agree with the oversimplification of their interface -- obfuscates too much. I wanted a
// structured logger which did not hide how to use it, was simple (i only want to implement a
// couple methods), and allowed all the complexity of what package was being used (e.g. zap,
// logrus, etc.) to be handled in a "New" function and a single call function -- aka "Print".
type Logger interface { 
	Print(level string, msg string, data ...map[string]interface{}) error
}