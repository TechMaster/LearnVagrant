
# build Goland binary to AMD architect to deploy
FROM golang:alpine AS build-env

WORKDIR /app

COPY go.mod .

COPY main.go .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w' -o main

RUN ls /app

# final stage
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy result binary go app to /app folder
COPY --from=build-env /app/main /app


ENTRYPOINT ["./main"]

EXPOSE 8001
