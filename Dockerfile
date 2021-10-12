FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN ls -l && \
    CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM alpine:latest
COPY --from=builder /app/app /app/app
WORKDIR /app
ENTRYPOINT [ "./app", "consumer" ]