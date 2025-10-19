# --- Stage 1: The Builder ---
# Use the official Go image as a base for building.
FROM golang:1.25.1-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the Go application.
# CGO_ENABLED=0 is important for creating a static binary.
# -o /app/main builds the output into a file named 'main'.
RUN CGO_ENABLED=0 go build -o /app/main ./cmd/api

# --- Stage 2: The Final Image ---
# Use a minimal, non-root image for the final container.
FROM gcr.io/distroless/static-debian11

# Set the working directory.
WORKDIR /app

# Copy only the compiled binary from the 'builder' stage.
COPY --from=builder /app/main .

# The command to run when the container starts.
CMD ["/app/main"]