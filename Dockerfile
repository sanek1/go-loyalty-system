FROM golang:1.23-alpine as builder
WORKDIR /app
RUN apk --no-cache add bash git make gcc gettext musl-dev

# depedency
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY . .
RUN go build -o ./bin/app cmd/gophermart/main.go

FROM alpine as runner

COPY --from=builder /app/bin/app /app
COPY  config/config.yaml config.yaml

CMD ["./app"]
