FROM gcr.io/voonik-tech/vnkshopup/golang-1-15-8:1.0 as base
# We use multi-stage build to avoid storing GitHub token inside Docker image
FROM base as builder
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=github.com/voonik

WORKDIR /ss2

# GitHub token for downloading private dependencies
ARG GITHUB_TOKEN

RUN git config --global url."https://${GITHUB_TOKEN}@github.com/voonik/".insteadOf "https://github.com/voonik/"

# Copy the go.mod and go.sum files to WORKDIR
COPY go.mod go.sum ./

# Download the dependencies from go.mod go.sum without checking complete code
RUN go mod download

COPY . .

# Install the dependencies without checking for go code
RUN go mod vendor

# Build the application
RUN go build -o ./ ./cmd/server ./cmd/worker

#Final Docker Image
FROM base
WORKDIR /ss2
COPY --from=builder /ss2 .
EXPOSE 5006
CMD ["./server"]
