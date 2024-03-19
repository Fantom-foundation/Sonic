package logger

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/exp/slog"
)

// init with defaults.
func init() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true)))
}

// SetDSN appends sentry hook to log root handler.
func SetDSN(value string) {
	panic("not implemented")
	/*
		// If DSN is empty, we don't create new hook.
		// Otherwise we'll the same error message for each new log.
		if value == "" {
			log.Warn("Sentry client DSN is empty")
			return
		}

		// TODO: find or make sentry log.Handler without logrus.
		sentry, err := logrus_sentry.NewSentryHook(value, nil)
		if err != nil {
			log.Warn("Probably Sentry host is not running", "err", err)
			return
		}

		log.Root().SetHandler(
			log.MultiHandler(
				log.Root().GetHandler(),
				LogrusHandler(sentry),
			))
	*/
}

// SetLevel sets level filter on log root handler.
// So it should be called last.
func SetLevel(l string) {
	lvl, err := levelFromString(l)
	if err != nil {
		panic(err)
	}

	log.SetDefault(
		log.NewLogger(
			log.NewTerminalHandlerWithLevel(os.Stderr, lvl, true),
		))
}

func levelFromString(lvlString string) (slog.Level, error) {
	switch lvlString {
	case "debug", "dbug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error", "eror":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unknown level: %v", lvlString)
	}
}
