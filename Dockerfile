FROM golang:1.22.6-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o app

FROM alpine:3.19
LABEL maintainer="K8sCat <k8scat@gmail.com>"
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
ENV GIN_MODE=release
ENTRYPOINT ["./app"]
