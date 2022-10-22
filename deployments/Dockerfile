FROM golang:1.16 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o credit_holidays_app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
RUN apk update \
  && apk -U upgrade \
  && apk add --no-cache ca-certificates bash gcc \
  && update-ca-certificates --fresh \
  && rm -rf /var/cache/apk/*
COPY --from=builder /app/credit_holidays_app /app/credit_holidays_app
COPY ./configs/config.yml /app/configs/config.yml
RUN chmod +x /app/credit_holidays_app
VOLUME /app/static
EXPOSE 8080
ENTRYPOINT ["/app/credit_holidays_app"]
