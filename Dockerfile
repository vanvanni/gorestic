FROM golang:1.23-alpine AS builder

RUN apk add -U tzdata
RUN apk --update add ca-certificates
RUN addgroup -S gorestic && adduser -S gorestic -G gorestic

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-w -s" -o gorestic

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/gorestic /gorestic

USER gorestic

CMD ["/gorestic"]