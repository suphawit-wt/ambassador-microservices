# Preparing build images
FROM golang:1.21.4-alpine3.18 AS builder
RUN apk add alpine-sdk

# Preparing dependencies and Build Application
WORKDIR /builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -tags musl -ldflags '-s' -o build

# Setup Runner Image for Application
FROM alpine:3.18.4 AS runner
WORKDIR /app
COPY --from=builder /builder/build .

# Execute Application
ENTRYPOINT [ "/app/build" ]