# =============================================================================
# Stage 1: Builder
# Uses the full Go toolchain image to compile the binary.
# This image is large (~800 MB) but is discarded after the build — it never
# ends up in the final image that gets deployed.
# =============================================================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy module files first and download dependencies.
# Docker caches each layer. Because go.mod/go.sum change far less often than
# source code, this layer is reused on subsequent builds — making them fast.
COPY go.mod go.sum ./
RUN go mod download

# Now copy source code. Changes here invalidate the cache from this point on,
# but the dependency download layer above is still reused.
COPY . .

# Build flags:
#   CGO_ENABLED=0  — disables C interop, producing a fully static binary
#   GOOS=linux     — cross-compile for Linux (required when building on macOS)
#   -ldflags="-w -s" — strips debug info, reducing binary size by ~30%
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o football-league ./cmd/api

# =============================================================================
# Stage 2: Runtime
# A minimal Alpine image (~7 MB). Contains only the compiled binary.
# No Go toolchain, no source code, no build artifacts.
# =============================================================================
FROM alpine:3.23

# ca-certificates is needed for any outbound HTTPS calls (not used now,
# but good practice so the image is not broken if networking is added later).
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the binary from the builder stage. Nothing else.
COPY --from=builder /app/football-league .

EXPOSE 8080

CMD ["./football-league"]