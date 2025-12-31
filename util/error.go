package util

import (
	"fmt"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

func LogDeferWarn(f func() error) {
	err := f()
	if err != nil {
		log.WithError(err).Warn("error in deferred job")
	}
}

func PanicToError(pan any, origErr error) error {
	if pan == nil {
		return origErr
	}
	log.WithField("stack", debug.Stack()).Info("trace of panic which occurred")
	switch x := pan.(type) {
	case string:
		return fmt.Errorf("panicked: %s", x)
	case error:
		return fmt.Errorf("panicked: %w", x)
	default:
		return fmt.Errorf("panicked: %v", x)
	}
}
