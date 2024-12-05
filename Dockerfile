FROM golang:1.23-alpine AS builder

RUN apk add -U tzdata ncurses-terminfo-base
RUN apk --update add ca-certificates

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-w -s" -o gorestic

RUN addgroup -S app && \
    adduser -S -D -H -h /home/app -s /sbin/nologin -G app -u 1000 app && \
    mkdir -p /home/app/.config && \
    chown -R app:app /home/app && \
    chmod -R 755 /home/app

RUN mkdir -p /home/app/.config && touch /home/app/.config/.keep

FROM scratch

WORKDIR /

COPY --from=builder /bin/sh /bin/sh
COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder /lib/libc.musl-x86_64.so.1 /lib/

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

COPY --from=builder /app/gorestic /gorestic

COPY --from=builder /home/app /home/app

VOLUME /home/app/.config

ENV GORESTIC_DOCKER=true
ENV HOME=/home/app
ENV UID=1000

USER app

CMD ["/gorestic"]