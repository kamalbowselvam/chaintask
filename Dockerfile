FROM golang:1.20.5-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app
COPY . .

RUN apk update
RUN apk add curl

# Build the Go app
RUN go build -o  main main.go
# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates
WORKDIR /app

COPY config /app/config
COPY --from=build_base /app/main .
COPY app.env .
RUN sed -i 's/localhost/postgres/g' app.env
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]