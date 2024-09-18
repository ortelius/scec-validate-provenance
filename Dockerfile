FROM cgr.dev/chainguard/go@sha256:ba65a2ad761366e9eebc7726d1dd1bc108fc2f6941cc86cc735b7a17451b014d AS builder

WORKDIR /app
COPY . /app

RUN go mod tidy; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:5915b6b2c9ff77916fc534d9c0676eaaf776964f674e4a6aeb6b43426c7db79a

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
