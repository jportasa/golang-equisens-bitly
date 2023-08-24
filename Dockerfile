ARG GOVERSION="1.19"

############################
# STEP 1 build executable binary
############################
FROM golang:${GOVERSION}-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
# Download Go modules
RUN go mod download
COPY *.go ./
# Build
RUN CGO_ENABLED=0 GOOS=linux ARCH=amd64 go build -o /main

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

############################
# STEP 2 create a small image with the binary
############################

FROM scratch
WORKDIR /
COPY --from=builder /main /main
USER nonroot:nonroot
CMD ["/main"]