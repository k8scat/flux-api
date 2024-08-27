FROM golang:1.22.6-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app

FROM golang:1.22.6-alpine
LABEL maintainer="K8sCat <k8scat@gmail.com>"
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
ENV GIN_MODE=release
ENTRYPOINT ["./app"]
