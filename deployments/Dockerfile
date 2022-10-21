FROM golang:1.16
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o credit_holidays_app ./cmd/main.go
CMD /app/credit_holidays_app