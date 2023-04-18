package statistics

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strconv"

	"github.com/peano88/fizzbazz-rest/pkg/model"
	"github.com/peano88/fizzbazz-rest/pkg/utils"
	"github.com/redis/go-redis/v9"
)

const (
	fizzBuzzStatisticsSet           = "fizzbuzz:statistics"
	redisDBAddressEnvVar            = "REDIS_DB_ADDRESS"
	redisDBTLSEnvVar                = "REDIS_DB_TLS"
	redisDBTLSInsecureEnvVar        = "REDIS_DB_TLS_INSECURE"
	redisDBUsernameEnvVar           = "REDIS_DB_USERNAME"
	redisDBPasswordEnvVar           = "REDIS_DB_PASSWORD"
	redisDBIdEnvVar                 = "REDIS_DB_ID"
	redisDBTLSCertificatePathEnvVar = "REDIS_DB_TLS_CERTIFICATE_PATH"
	redisDBTLSKeyPathEnvVar         = "REDIS_DB_TLS_KEY_PATH"
)

// FizzBuzzStatsRedis is a statistic component based on redis DB
type FizzBuzzStatsRedis struct {
	rdb *redis.Client
}

// NewFizzBuzzStatsRedis instances a new FizzBuzzStatsRedis, which will automatically handles reconnection
// to the redis DB. the address can be set using environment variable REDIS_DB_ADDRESS as well as user login 
// via variables REDIS_DB_USERNAME, REDIS_DB_PASSWORD. The db instance to use is not provide by default and can
// be changed via REDIS_DB_ID. If the connection needs a TLS protection, than variable REDIS_DB_TLS needs to be set to 
// true. When using TLS the configuration can be tweaked using REDIS_DB_TLS_INSECURE, REDIS_DB_TLS_CERTIFICATE_PATH,
// REDIS_DB_TLS_KEY_PATH which will allow for an insecure connection and specific client certificate+key.
// Standard variable SSL_CERT_FILE and SSL_CERT_DIR can be used to change the default loading of system CAs.
func NewFizzBuzzStatsRedis() (*FizzBuzzStatsRedis, error) {
	redisOptions := redis.Options{
		Addr:     utils.GetEnv(redisDBAddressEnvVar, "localhost:6379"),
		Username: utils.GetEnv(redisDBUsernameEnvVar, ""),
		Password: utils.GetEnv(redisDBPasswordEnvVar, ""),
	}

	if id, err := strconv.Atoi(utils.GetEnv(redisDBIdEnvVar, "0")); err == nil {
		redisOptions.DB = id
	}

	if utils.IsTLSEnabled(redisDBTLSEnvVar) {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("error loading system certpool: %w", err)
		}

		certificatePath := os.Getenv(redisDBTLSCertificatePathEnvVar)
		keyPath := os.Getenv(redisDBTLSKeyPathEnvVar)
		tlsCertificates := []tls.Certificate{}

		if certificatePath != "" && keyPath != "" {
			certificate, err := tls.LoadX509KeyPair(certificatePath, keyPath)
			if err != nil {
				return nil, fmt.Errorf("error loading tls certificate: %w", err)
			}
			tlsCertificates = append(tlsCertificates, certificate)
		}

		redisOptions.TLSConfig = &tls.Config{
			Certificates:       tlsCertificates,
			RootCAs:            certPool,
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}

	}

	return &FizzBuzzStatsRedis{
		rdb: redis.NewClient(&redisOptions),
	}, nil
}

// Increment uses redis ZINCRBY to increment the request count of the provided set of input parameters. The set identifier
// is built by concatenation of each parameter using the model.Separator
func (fs *FizzBuzzStatsRedis) Increment(ctx context.Context, n, m, top int, fizz, buzz string) error {

	member := strconv.Itoa(n) + model.Separator + strconv.Itoa(m) + model.Separator + strconv.Itoa(top) + model.Separator + fizz + model.Separator + buzz

	if err := fs.rdb.ZIncrBy(ctx, fizzBuzzStatisticsSet, 1.0, member).Err(); err != nil {
		return fmt.Errorf("error in incrementing input parameters counter: %w", err)
	}

	return nil
}

// Stats will return the most request set using ZREVRANGEBYSCORE of redis.Will return NoStatsAvailable error
// if no statistic of previous requests is available
func (fs *FizzBuzzStatsRedis) Stats(ctx context.Context) (model.FizzBuzzStatisticsOutput, error) {
	res, err := fs.rdb.ZRevRangeByScoreWithScores(ctx, fizzBuzzStatisticsSet, &redis.ZRangeBy{
		Min:    "0",
		Max:    "+inf",
		Offset: 0,
		Count:  1,
	}).Result()

	if err != nil {
		return model.FizzBuzzStatisticsOutput{}, err
	}

	if len(res) == 0 {
		return model.FizzBuzzStatisticsOutput{}, NoStatsAvailable{}
	}

	return utils.FizzBuzzStatisticsOutputFromString(res[0].Member.(string), int64(res[0].Score))
}
