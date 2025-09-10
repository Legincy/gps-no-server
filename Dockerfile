FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gps-no-server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/gps-no-server .

RUN mkdir -p /app/data

#RUN go install github.com/swaggo/swag/cmd/swag@latest

#RUN swag init -g cmd/server/main.go -o ./.devcontainer/docs

# Funktioniert noch nicht. Bis das funktioniert muss bei änderung der Routen kommentare im Hauptordner "swag init -g cmd/server/main.go -o ./.devcontainer/docs" ausgeführt werden :)

EXPOSE 8080

CMD ["./gps-no-server"]