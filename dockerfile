########################################################
# STEP 1 build executable binary
########################################################
FROM golang:alpine AS builder

# Install git for fetching the dependencies.
RUN apk update && apk add --no-cache git curl 

WORKDIR $GOPATH/src/app/

COPY ./ ./

# Fetch dependencies using go get.
RUN go mod download

# Build the binary.
RUN GOOS=linux go build -ldflags="-w -s" -o /statement_processor ./main.go

########################################################
# STEP 2 build image
########################################################
FROM scratch

# Copy our static executable.
COPY --from=builder /statement_processor /app/statement_processor
WORKDIR /app/

# Run the binary.
ENTRYPOINT ["/app/statement_processor"]