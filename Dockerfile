FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Build 
RUN go build -o url-shortener .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/url-shortener .

#run
CMD ["./url-shortener"]
