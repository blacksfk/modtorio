FROM golang:alpine AS builder

WORKDIR /go/src/modtorio

# copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source files
COPY . ./

# compile
RUN go build

# separate image for running
FROM alpine

# create container working directory
WORKDIR /app/

# copy compiled executable
COPY --from=builder /go/src/modtorio/modtorio .

# start the tool
ENTRYPOINT ["./modtorio"]
