FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/mailhook .

FROM alpine:latest AS runtime
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=build /app/mailhook /app/mailhook
RUN chmod +x /app/mailhook
ENTRYPOINT ["/app/mailhook"]
EXPOSE 25