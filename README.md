# FizzBuzz REST server

<!--toc:start-->
- [FizzBuzz REST server](#fizzbuzz-rest-server)
  - [Implementation](#implementation)
  - [Configuration](#configuration)
  - [Documentation](#documentation)
  - [Redis](#redis)
  - [Dev corner](#dev-corner)
<!--toc:end-->

[FizzBuzz](https://en.wikipedia.org/wiki/Fizz_buzz) is a very common word game which consists of creating sequence of strings. 
This implementation via a REST api server allows to:
- replace 3 and 5 with custom values;
- replace `fizz` and `buzz` with custom values;
- set a custom value as inclusive upper limit of the sequence.

Furthermore, it is possible to receive which set of input parameters is the most requested.

## Implementation

The project consists of an HTTP api server built using Go. It exposes two endpoints under `/api/v1` on port `3000`:

1. `/fizzbuzz` (GET): returns a fizz-buzz-alike sequence based on the following query parameters (all required): `int1`, `int2`, `limit`, `str1`, `str2` where the
first twos will provides the two base numbers for the sequence (much like 3 and 5 in the original version), `limit` is the inclusive upper limit of the sequence and
`str1` and `str2` are the two strings corresponding to `fizz` and `buzz` in the original version. Optionally, query parameter `start` (defaulted to 1) can be used to 
start the sequence in a given position. If the requested sequence has more than 65536 elements, than only the first
65536 items are returned together with a link to a `fizzbuzz` request which will extend/complete the sequence.
2. `/statistics` (GET): return the set of query parameters, which corresponds to the most demanded request on GET `/fizzbuzz`. The hit-count (number of received request for
the set of parameters) is returned as well. If two (or more) sets share the same hit-count, then the sets are order by reversed lexicographical order and the first set is returned. 


The statistics part is implemented using a [redis DB](https://redis.io/). 

## Configuration

The application uses a [12-factor](https://12factor.net/) approach on the configuration and can be modified via environment variables. Here a brief description:
| Variable | Usage | Allowed values |
| --- | --- | --- |
| FIZZBUZZ_LOG_LEVEL | Set the log level of the application; defaulted to `info` | `panic`, `error`, `warn`, `info`, `debug`, `trace` | 
| FIZZBUZZ_TLS_ENABLE | The server will listen for TLS connection | same string values compatibles with go `strconv.ParseBool` |
| FIZZBUZZ_TLS_INSECURE | Allows insecure connection, defaulted to `false` | same string values compatibles with go `strconv.ParseBool` |
| FIZZBUZZ_CLIENT_AUTH_TYPE | Force provided client authentication type | same as `tls.Config` |
| FIZZBUZZ_TLS_CERT | Path of the server certificate for TLS. Mandatory if TLS is enabled | |
| FIZZBUZZ_TLS_KEY | Path of the server key for TLS. Mandatory if TLS is enabled | |


## Documentation

The REST api is documented in OpenAPI 3.0 format in the [openapi file](./openapi.yaml). 

The Go project documentation can be generated using:
```bash
$ godoc -http=:6060
```

and can be explored by pointing the browser to `localhost:6060`.

## Redis

The connection to the redis DB can be parametrized using these environment variables:

| Variable | Usage | Allowed values |
| --- | --- | --- |
| REDIS_DB_ADDRESS | addres of the redis instance | |
| REDIS_DB_USERNAME | username | |
| REDIS_DB_PASsWORD | password| |
| REDIS_DB_ID | numeric id of the DB, defaulted to 0 | integer |
| REDIS_DB_TLS | use TLS to establish the connection | same string values compatibles with go `strconv.ParseBool` |
| REDIS_DB_TLS_INSECURE | allows for insecure connection | same string values compatibles with go `strconv.ParseBool` |
| REDIS_DB_TLS_CERTIFICATE_PATH | path of the client certificate | |
| REDIS_DB_TLS_KEY_PATH | path of the client key | |

if the redis DB is not available then the GET `/fizzbuzz` will still answer as expected (altough the incoming requests are not automatically registered for statistics pupose) whereas the GET `/statistics` will return a `503 Service unavailable` response.
The application handles the reconnection automatically.

## Dev corner
Use [nix](https://nixos.org/) to create the development environment. A file [shell.nix](./shell.nix) is available at the root of the repository.

---

The easiest way to get a working redis DB instance is using docker:
```bash
$ docker run --name my-redis -p 6379:6379 -d redis:alpine
```
which can be stopped/restarted via
```bash
$ docker stop my-redis
$ docker start my-redis
```

---

Unit tests for `pkg/server` requires mocking of interface `server.FizzBuzzStats`. The mock has been generated using [vektra/mockery](https://github.com/vektra/mockery).
Please refer to the official documentation for any issue.


