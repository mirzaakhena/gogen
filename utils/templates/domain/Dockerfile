################
# BUILD BINARY #
################

FROM golang:1.19.4-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=v1.0.0'" .

#####################
# MAKE SMALL BINARY #
#####################
FROM scratch

# Copy the executable.
WORKDIR /app

COPY --from=builder /app/{{.PackagePath}} /usr/bin/
COPY --from=builder /app/config.json /app

# ENTRYPOINT ["{{.PackagePath}}", "myapp"]