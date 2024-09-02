FROM cgr.dev/chainguard/go@sha256:b709bdeca91b88b895bfe08b92776ee4db559fe94e30044ae7244ad5aa4c785d AS builder

WORKDIR /app
COPY . /app

RUN go mod tidy; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:09c74721beb055d94e49515381cc5118d85bae208e7ddefe1d28f6b355a1a6f5

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs docs

ENV ARANGO_HOST localhost
ENV ARANGO_USER root
ENV ARANGO_PASS rootpassword
ENV ARANGO_PORT 8529
ENV MS_PORT 8080

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]
