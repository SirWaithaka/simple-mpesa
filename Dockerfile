# Using multi-stage builds
FROM golang:1.14-alpine AS go-builder

LABEL maintainer="Sir Waithaka"

# Set the current working directory inside the container
WORKDIR /go/src/application

# copy go.{mod,sum} files for use to fetch dependencies
# fetching go dependencies first allows the build tool to cache this part of the image
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy project source files
COPY app/ ./app
COPY cmd/ ./cmd
COPY configs/ ./configs

# Build the application
RUN mkdir bin/
RUN go build -o bin/wallet-server cmd/wallet-server.go


# Start the second image
FROM alpine:3

# install some linux packages
RUN apk --no-cache add ca-certificates tzdata

# configure correct timezone for image
RUN cp /usr/share/zoneinfo/Africa/Nairobi /etc/localtime
RUN echo "Africa/Nairobi" > /etc/timezone

# set the working director in the container
WORKDIR /root/app

# copy extra files that can be useful to someone reading the application image
COPY Dockerfile/ ./Dockerfile
COPY ReadMe.md/ ./ReadMe.md

# copy application configuration from example file
COPY config.yml.example config.yml

COPY --from=go-builder /go/src/application/bin/wallet-server .

# expose the port that the server starts on
EXPOSE 6700

CMD ["./wallet-server"]
