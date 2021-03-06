ARG GO_VERSION=1.16.6

FROM golang:${GO_VERSION}-alpine AS builder

# Get dependencies directly
RUN go env -w GOPROXY=direct
# It's necessary to install git when building the app
RUN apk add --no-cache git
# Security certificates
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src
# Copy both files to the main directory
COPY ./go.mod ./go.sum ./
# Download and install dependencies
RUN go mod download

COPY ./ ./
# Disable the C++ compiler
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    # Select the name
    -o /api-ws

# This image will run the app
FROM scratch as runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY .env ./

COPY --from=builder /api-ws /api-ws

EXPOSE 5050
# Call the server
ENTRYPOINT ["/api-ws"]