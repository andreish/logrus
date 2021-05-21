package logrus

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestIssue1201(t *testing.T) {
	t.Log("TestIssue1201")
	setupLogging()
	c := make(chan int)
	go testLogging(c)
	last := -1
	timeout := false
	for k := 0; k < 10 && (!timeout); k++ {
		select {
		case last = <-c:
			continue
		case <-time.After(30 * time.Second): //30 secs shall be enough
			t.Log(fmt.Sprintf("timeout last=%d", last))
			timeout = true
		}
	}
	if last != 9 {
		t.Errorf("got last =%d wanted 9", last)
	}
	t.Log(fmt.Sprintf("got last = %d ", last))
}

type fileLogWriter struct {
	filename string
}

func newFileLogWriter(filename string) *fileLogWriter {
	return &fileLogWriter{filename: filename}
}

func (fw *fileLogWriter) Write(p []byte) (n int, err error) {
	Info(fmt.Sprintf("opening %s", fw.filename))
	logfile, err := os.OpenFile(fw.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Error(fmt.Sprintf("Failed to open logfile <%s>", fw.filename), err)
		return
	}
	defer logfile.Close()
	return logfile.Write(p)
}

var logger *Entry

func setupLogging() *Entry {

	SetFormatter(&TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	logWriter := newFileLogWriter("./zzztest.log")

	logoutput := io.MultiWriter(os.Stdout, logWriter)
	SetOutput(logoutput)
	SetLevel(InfoLevel)
	logger = WithFields(Fields{
		"service": "qa",
	})
	return logger
}

func testLogging(ch chan<- int) {
	for i := 0; i < 10; i++ {
		logger.Info(fmt.Sprintf("log line %d ", i))
		ch <- i
	}
	close(ch)
}
