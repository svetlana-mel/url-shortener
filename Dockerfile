FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git gcc gettext musl-dev

# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

# build
COPY ./app ./
RUN go build -o ./bin/url-shortener cmd/url-shortener/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/url-shortener /

WORKDIR /

COPY ./app/config/prod.yaml /config/prod.yaml
# COPY ./app/config/local.yaml /config/local.yaml
COPY ./app/.env /.env

# setup database
RUN apk add --no-cache sqlite
RUN mkdir -p /storage && \
    sqlite3 /storage/storage.db "VACUUM;" && \
    chmod 644 /storage/storage.db

EXPOSE 8082

CMD [ "/url-shortener" ]