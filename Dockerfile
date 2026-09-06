FROM golang:1.25-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata wget \
 && adduser -D -H -u 10001 app

WORKDIR /app

COPY --from=build /out/api /app/api
COPY --from=build /out/seed /app/seed

USER app

EXPOSE 8080

ENTRYPOINT ["/app/api"]
