package log

import (
	"fmt"
	"strings"
	"time"
)

// INFO, DEBUG, and ERROR are the 3 levels for this
const (
	INFO  = 0
	DEBUG = 1
	ERROR = 2
)

type message struct {
	level     int
	timestamp time.Time
	text      string
	params    []interface{}
}

var queue chan message
var done chan struct{}
var level int

// Start opens the log queue
func Start(lvl int) {
	queue = make(chan message, 50)
	done = make(chan struct{})
	level = lvl

	Msg(0, "Starting log service")
	go func() {
		for msg := range queue {
			printMsg(msg)
		}
		done <- struct{}{}
	}()
}

// Close stops the running service
func Close() {
	Msg(0, "log service closed")
	close(queue)
	<-done
}

func printMsg(msg message) {
	if msg.level > level {
		return
	}
	text := msg.text
	if len(msg.params) > 0 {
		text = fmt.Sprintf(text, msg.params...)
	}

	fullMsg := fmt.Sprintf("[%v] %v", msg.timestamp, text)
	fmt.Println(fullMsg)
}

// Msg queues up a message
func Msg(lvl int, msg string) {
	Msgf(lvl, msg)
}

// Msgf queues up a formatted message
func Msgf(lvl int, msg string, params ...interface{}) {
	ts := time.Now()

	queue <- message{
		level:     lvl,
		timestamp: ts,
		text:      strings.ToLower(msg),
		params:    params,
	}
}
