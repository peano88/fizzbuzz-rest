package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/peano88/fizzbazz-rest/pkg/model"
	"github.com/peano88/fizzbazz-rest/pkg/utils"
)

const (
	TLSEnvVar            = "FIZZBUZZ_TLS_ENABLE"
	insecureEnvVar       = "FIZZBUZZ_INSECURE"
	clientAuthTypeEnvVar = "FIZZBUZZ_CLIENT_AUTH_TYPE"
	logLevelEnvVar       = "FIZZBUZZ_LOG_LEVEL"
)

//go:generate mockery --name FizzBuzzStats
//FizzBuzzStats is the interface representing what is expected by the statistics component
type FizzBuzzStats interface {
    // Increment receives the input parameters so that they can be registered
	Increment(ctx context.Context, n, m, top int, fizz, buzz string) error
    // Stats should return the model.FizzBuzzStatisticsOutput representing the #1 hit for the GET /fizzbuzz
    // error otherwise
	Stats(ctx context.Context) (model.FizzBuzzStatisticsOutput, error)
}

// FizzBuzzServer is the structure defining the HTTP requests handling and middleware
type FizzBuzzServer struct {
    // instance of FizzBuzzStats
	Stats FizzBuzzStats
}

// Configure will return a configured *http.Server which can be used to serve requests
// if environment variable FIZZBUZZ_TLS_ENABLE is set to true, than the server will be configured to be used
// with ListenAndServeTLS method. in this case, the TLS configuration will allow insecure connection when 
// FIZZBUZZ_TLS_INSECURE is set to true; the client authentification type can be changed using environment variable
// FIZZBUZZ_CLIENT_AUTH_TYPE
// Standard variable SSL_CERT_FILE and SSL_CERT_DIR can be used to change the default loading of system CAs.
// The serve will create a unique identifier for each incoming request, will log each request processing based on 
// variable FIZZBUZZ_LOG_LEVEL and will automatically recover from panics
func (fbs *FizzBuzzServer) Configure() (*http.Server, error) {
	logger := httplog.NewLogger("fizzbuzz-rest", httplog.Options{
		LogLevel: utils.GetEnv(logLevelEnvVar, "info"),
		JSON:     true,
	})

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Route("/fizzbuzz", func(r chi.Router) {
		r.Use(fbs.ValidationMiddleware)
		r.Use(fbs.ToStatisticsMiddleware)
		r.Get("/", fbs.GetFizzBuzzHandler)
	})

	r.Get("/statistics", fbs.GetStatisticsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Mount("/api/v1/", r)

	s := http.Server{
		Addr:         ":3000",
		Handler:      apiRouter,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	if utils.IsTLSEnabled(TLSEnvVar) {

		insecure, err := strconv.ParseBool(utils.GetEnv(insecureEnvVar, "false"))
		if err != nil {
			return nil, fmt.Errorf("error with insecure tls env var: %w", err)
		}

		clientAuthType, err := strconv.Atoi(utils.GetEnv(clientAuthTypeEnvVar, strconv.Itoa(int(tls.RequireAndVerifyClientCert))))
		if err != nil {
			return nil, fmt.Errorf("error with client auth type env var: %w", err)
		}

		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("error loading system cert pool: %w", err)
		}

		s.TLSConfig = &tls.Config{
			ClientAuth:         tls.ClientAuthType(clientAuthType),
			ClientCAs:          certPool,
			InsecureSkipVerify: insecure,
			CipherSuites:       []uint16{},
			MinVersion:         tls.VersionTLS12,
		}
	}

	return &s, nil
}
