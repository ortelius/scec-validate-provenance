FROM cgr.dev/chainguard/go@sha256:598027417e0a039dc326c958feb5a088447bec198ad74207d854558106b3318f AS builder

WORKDIR /app
COPY . /app

RUN go install github.com/swaggo/swag/cmd/swag@latest; \
    swag init; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:52dcdf7ebd60c68a59c0b635fea69a819a43c4cabf080958f0ec58326a64c772

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
