package testlog

import (
	"log"
	"os"
)

type testLogger struct {
	testName string
	logger   *log.Logger
}

func (l *testLogger) Log(message string) {
	l.logger.Printf("[%s] %s", l.testName, message)
}

func (l *testLogger) SetTestName(name string) {
	l.testName = name
}

var Logger *testLogger

// init function
func init() {
	Logger = &testLogger{
		testName: "default-test",
		logger:   log.New(os.Stdout, "[TERRATEST] ", log.LstdFlags),
	}
}


