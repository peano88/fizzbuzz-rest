package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/peano88/fizzbuzz-rest/pkg/server"
	"github.com/peano88/fizzbuzz-rest/pkg/statistics"
	"github.com/peano88/fizzbuzz-rest/pkg/utils"
)

const (
	certPathEnvVar = "FIZZBUZZ_TLS_CERT"
	keyPathEnvVar  = "FIZZBUZZ_TLS_KEY"
)

func main() {
	ctx, cancelMain := signal.NotifyContext(context.TODO(), os.Interrupt, os.Kill)
	defer cancelMain()

	fizzBuzzStats, err := statistics.NewFizzBuzzStatsRedis()
	if err != nil {
		log.Fatalf("error instantiating fizzbuzz statistics component: %s", err.Error())
	}

	fizzbuzzServer := server.FizzBuzzServer{
		Stats: fizzBuzzStats,
	}

	s, err := fizzbuzzServer.Configure()
	if err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error)

	if utils.IsTLSEnabled(server.TLSEnvVar) {
		go func() {
			if err := s.ListenAndServeTLS(os.Getenv(certPathEnvVar), os.Getenv(keyPathEnvVar)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}()

	} else {
		go func() {
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
		// graceful shutdown
		ctxShutdown, cancelShutdown := context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancelShutdown()
		if err := s.Shutdown(ctxShutdown); err != nil {
			log.Fatalf("could not shutdown gracefully: %s", err.Error())
		}
	case err := <-errChan:
		log.Fatalf("fatal error: %s", err.Error())
	}

}
