FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /validator-api .

FROM alpine:3.20
RUN adduser -D -H appuser
COPY --from=builder /validator-api /validator-api
USER appuser
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/validator-api"]
