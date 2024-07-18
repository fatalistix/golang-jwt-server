package worker

import (
	"log/slog"
	"time"

	"github.com/fatalistix/slogattr"
)

type ExpiredTokenDeleter interface {
	DeleteExpired() error
}

type RefreshTokenCleaner struct {
	log      *slog.Logger
	stopChan chan bool
	timeout  time.Duration
	deleter  ExpiredTokenDeleter
}

func MakeRefreshTokenCleaner(log *slog.Logger, timeout time.Duration, d ExpiredTokenDeleter) RefreshTokenCleaner {
	return RefreshTokenCleaner{
		log:      log,
		stopChan: make(chan bool, 1),
		timeout:  timeout,
		deleter:  d,
	}
}

func (w RefreshTokenCleaner) Start() {
	go func() {
		const op = "worker.RefreshTokenCleaner.Start@lambda"

		log := w.log.With(
			slog.String("op", op),
		)

		ticker := time.NewTicker(w.timeout)
		defer ticker.Stop()
		for {
			select {

			case <-ticker.C:
				err := w.deleter.DeleteExpired()
				if err != nil {
					log.Error("worker stopped with error", slogattr.Err(err))
					return
				}
				log.Info("expired tokens cleaned")

			case <-w.stopChan:
				log.Info("worker stopped")
				return

			}
		}
	}()
}

func (w RefreshTokenCleaner) Stop() {
	w.stopChan <- true
}
