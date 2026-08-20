# ---------- Build stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache module downloads separately from source changes
COPY go.mod ./
RUN go mod download

COPY . .

# Static binary, no cgo, so it runs on the minimal runtime image below
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/portfolio .

# ---------- Runtime stage ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/portfolio ./portfolio
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

USER app

ENV PORT=8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/ || exit 1

ENTRYPOINT ["./portfolio"]
