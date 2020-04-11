package log

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/op/go-logging"
)

var loggingEnable = true
var log = logging.MustGetLogger("Logger")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{time:2006-01-02 15:04:05.999} [%{level:.1s}] %{message}`,
)

// Secure is an an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
type Secure string

func init() {
	SetupLogger(true, "DEBUG")
}

// Redacted is called whenever anything is logged using Secure
func (p Secure) Redacted() interface{} {
	return logging.Redact(string(p))
}

// SetupLogger is called in initialization part of this service
func SetupLogger(isEnabled bool, level string) {
	loggingEnable = isEnabled
	backend := logging.NewLogBackend(os.Stdout, "", 0)

	// For messages written to backend we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backendFormatter := logging.NewBackendFormatter(backend, format)

	var lvl logging.Level

	switch level {
	case "ERROR":
		lvl = logging.ERROR
	case "WARNING":
		lvl = logging.WARNING
	case "INFO":
		lvl = logging.INFO
	case "DEBUG":
		lvl = logging.DEBUG
	default:
		lvl = logging.INFO
	}

	// Only errors and more severe messages should be sent to backend1
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(lvl, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)
}

// D writes debug level log
func D(format string, v ...interface{}) {
	if loggingEnable {
		log.Debugf(format, v...)
	}
}

// W writes warning level log
func W(format string, v ...interface{}) {
	if loggingEnable {
		log.Warningf(format, v...)
	}
}

// E writes error level log
func E(format string, v ...interface{}) {
	if loggingEnable {
		log.Errorf(format, v...)
	}
}

// I writes info level log
func I(format string, v ...interface{}) {
	if loggingEnable {
		log.Infof(format, v...)
	}
}

// GetRequestID generate new ID for request.
func GetRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
