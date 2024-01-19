# build backend
FROM golang:alpine as go-builder
WORKDIR /app
COPY go.* ./
RUN go mod download
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./kahva ./cmd

# build frontend
FROM node:lts-alpine as vue-builder
WORKDIR /app
COPY web/package*.json ./
RUN npm install
COPY web/. .
RUN npm run build

FROM scratch
# copy certs if remote xmlrpc server is using TLS
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# copy backend binary
COPY --from=go-builder /app/kahva /app/kahva
# copy frontend
COPY --from=vue-builder /app/dist /app/www
EXPOSE 8080
WORKDIR /app
CMD ["/app/kahva"]