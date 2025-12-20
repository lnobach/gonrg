package util

import log "github.com/sirupsen/logrus"

func LogDeferWarn(f func() error) {
	err := f()
	if err != nil {
		log.WithError(err).Warn("error in deferred job")
	}
}
