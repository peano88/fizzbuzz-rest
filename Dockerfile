FROM golang:1.20-alpine AS builder
COPY $PWD /fizzbuzz-rest
RUN cd /fizzbuzz-rest \
    && go build cmd/fizzbuzz/main.go

FROM alpine
LABEL maintainer="peano88 <ilpeano@gmail.com>"

COPY --from=builder /fizzbuzz-rest/main /usr/local/fizzbuzz/bin/fizzbuzz
ENTRYPOINT [ "/usr/local/fizzbuzz/bin/fizzbuzz" ]


