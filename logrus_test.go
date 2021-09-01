package coolog

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestLogrus(t *testing.T) {
	fmt.Println("Starting TestLogrus...")

	// Taken from https://stackoverflow.com/a/29339052
	// We need to capture the standard out so when we log to it, we can verify it works.

	// Backup the stdout so we can return it after we are finished.
	backupStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Start the logrus logger.
	l, err := NewLogrusLogger("info", "console", "text")
	if err != nil {
		t.Errorf("While creating NewLogrusLogger, expect err: nil, got err: %v", err)
		return
	}

	// Testing begins...
	infoMsg := "I am info, hear me roar!"
	infoExtrasName := "infoExtras"
	infoExtras := "extras..."

	l.Print("info", infoMsg, map[string]interface{}{infoExtrasName: infoExtras})

	errorMsg := "Ohhhhhh Noooooooooo Mr. Bill!"

	errorExtras1Name := "errorExtras1"
	errorExtras1 := "first thing"
	errorExtras2Name := "errorExtras2"
	errorExtras2 := 666
	errorExtras3Name := "errorExtras3"
	errorExtras3 := "poop"

	l.Print("error", errorMsg, map[string]interface{}{
		errorExtras1Name: errorExtras1,
		errorExtras2Name: errorExtras2,
		errorExtras3Name: errorExtras3,
	})

	l.Print("trace", "I should be invisible!", map[string]interface{}{"wha": "As should I..."})

	w.Close()
	//l.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = backupStdout

	// Checking begins...
	logs := strings.Split(string(out), "\n")

	if len(logs) < 1 {
		t.Error("While checking the log messages, expected length of logs to be 2, got 0")
		return
	}

	if !strings.Contains(logs[0], infoMsg) ||
		!strings.Contains(logs[0], infoExtrasName+"="+infoExtras) {
		t.Errorf("While checking the log messages, expected proper info message, got %s", logs[0])
		return
	}

	if len(logs) < 2 {
		t.Error("While checking the log messages, expected length of logs to be 2, got 1")
		return
	}

	if !strings.Contains(logs[1], errorMsg) ||
		!strings.Contains(logs[1], errorExtras1Name+"=\""+errorExtras1+"\"") ||
		!strings.Contains(logs[1], errorExtras2Name+"="+strconv.Itoa(errorExtras2)) ||
		!strings.Contains(logs[1], errorExtras3Name+"="+errorExtras3) {
		t.Errorf("While checking the log messages, expected proper error message, got %s", logs[1])
		return
	}

	if len(logs) > 2 && logs[2] != "" {
		t.Errorf("While checking the log messages, expected length of logs to be 2, got %d\t%+v", len(logs), logs)
		return
	}

	fmt.Println("TestLogrus\tOK")
}
