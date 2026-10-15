FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /release-api .

FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /release-api /usr/local/bin/release-api
USER app
EXPOSE 8080
ENV PORT=8080
CMD ["/usr/local/bin/release-api"]
