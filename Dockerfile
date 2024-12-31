FROM cgr.dev/chainguard/go@sha256:290786b941bf42bf9913948dc008086877f67d646f5641590d878d930f458803 AS builder

WORKDIR /app
COPY . /app

RUN go mod tidy; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:1e61f5abe1233356c3e63ab4e050b846cc0ffcb5e6516da8e74cd81982bd9328

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
