FROM golang:1.23-alpine AS builder

RUN apk add -U tzdata
RUN apk --update add ca-certificates

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-w -s" -o gorestic

RUN mkdir -p /home/root/.config && touch /home/root/.config/.keep

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

COPY --from=builder /app/gorestic /gorestic

COPY --from=builder /home/root /home/root

VOLUME /home/gorestic/.config

ENV HOME=/home/root

CMD ["/gorestic"]