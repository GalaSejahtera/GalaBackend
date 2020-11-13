FROM golang:1.13-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd ./cmd/
COPY pkg ./pkg/

RUN go build -o ./out/main ./cmd/server/main.go

FROM alpine:3.9
RUN apk add ca-certificates tzdata

ENV TEMPLATES_PATH=/app/templates
COPY templates /app/templates/

COPY --from=build_base /tmp/app/out/main /app/main

EXPOSE 10001
EXPOSE 10002

ENTRYPOINT ["/app/main"]