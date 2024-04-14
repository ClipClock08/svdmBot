# Use the official Go image as the base image
FROM golang:latest AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o main .

# Use a lightweight image for the final build
FROM gcr.io/distroless/base-debian10

# Copy the built executable from the previous stage
COPY --from=build /app/main /

# Set the entry point for the container
ENTRYPOINT ["/main"]