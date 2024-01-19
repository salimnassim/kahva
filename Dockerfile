FROM golang:alpine as go-builder
WORKDIR /app
COPY go.* ./
RUN go mod download
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./kahva ./cmd

FROM scratch
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /app/kahva /app/kahva
EXPOSE 8080
WORKDIR /app
CMD ["/app/kahva"]