FROM golang:1.16 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o credit_holidays_app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
RUN apk update \
  && apk -U upgrade \
  && apk add --no-cache ca-certificates bash gcc \
  && update-ca-certificates --fresh \
  && rm -rf /var/cache/apk/*
RUN addgroup egor && adduser -S egor -u 1000 -G egor
COPY --chown=egor:egor --from=builder /app/credit_holidays_app .
RUN chmod +x /app/credit_holidays_app
USER egor
EXPOSE 8080
CMD ["/app/credit_holidays_app"]
