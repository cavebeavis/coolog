package coolog

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestZap(t *testing.T) {
	fmt.Println("Starting TestZap...")

	// Taken from https://stackoverflow.com/a/29339052
	// We need to capture the standard out so when we log to it, we can verify it works.

	// Backup the stdout so we can return it after we are finished.
	backupStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Start the zap logger.
	z, err := NewZapLogger("info", []string{"console"}, "plain")
	if err != nil {
		t.Errorf("While creating NewZapLogger, expect err: nil, got err: %v", err)
		return
	}

	// Testing begins...
	infoMsg := "I am info, hear me roar!"
	infoExtrasName := "infoExtras"
	infoExtras := "extras..."

	z.Print("info", infoMsg, map[string]interface{}{infoExtrasName: infoExtras})

	errorMsg := "Ohhhhhh Noooooooooo Mr. Bill!"
	errorExtras1Name := "errorExtras1"
	errorExtras1 := "first thing"
	errorExtras2Name := "errorExtras2"
	errorExtras2 := 666
	errorExtras3Name := "errorExtras3"
	errorExtras3 := "poop"

	z.Print(
		"error",
		errorMsg,
		map[string]interface{}{
			errorExtras1Name: errorExtras1,
			errorExtras2Name: errorExtras2,
			errorExtras3Name: errorExtras3,
		},
	)

	traceMsg := "I should be invisible!"
	z.Print("trace", traceMsg, map[string]interface{}{"wha": "As should I..."})

	w.Close()
	z.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = backupStdout

	// Checking begins...
	logs := strings.Split(string(out), "\n")

	if len(logs) < 1 {
		t.Error("While checking the log messages, expected length of logs to be 2, got 0")
		return
	}

	if !strings.Contains(logs[0], infoMsg) ||
		!strings.Contains(logs[0], "{\""+infoExtrasName+"\": \""+infoExtras+"\"}") {
		t.Errorf("While checking the log messages, expected proper info message, got %s", logs[0])
		return
	}

	if len(logs) < 2 {
		t.Error("While checking the log messages, expected length of logs to be 2, got 1")
		return
	}

	if !strings.Contains(logs[1], errorMsg) ||
		!strings.Contains(logs[1], "\""+errorExtras1Name+"\": \""+errorExtras1+"\"") ||
		!strings.Contains(logs[1], "\""+errorExtras2Name+"\": "+strconv.Itoa(errorExtras2)) ||
		!strings.Contains(logs[1], "\""+errorExtras3Name+"\": \""+errorExtras3+"\"") {
		t.Errorf("While checking the log messages, expected proper error message, got %s", logs[1])
		return
	}

	for _, v := range logs {
		if strings.Contains(v, traceMsg) {
			t.Errorf("While checking the log messages, expected no trace logs, got %#v", logs)
			return
		}
	}

	//t.Error("im a test spot to make sure this actually does shit...")

	fmt.Println("TestZap\tOK")
}
