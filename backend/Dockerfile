FROM golang:1.26-bookworm AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/worker ./cmd/worker

FROM debian:bookworm-slim AS api
RUN useradd --system --uid 10001 --create-home recode \
    && mkdir -p /data \
    && chown recode:recode /data
COPY --from=build /out/api /usr/local/bin/recode-api
USER recode
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/recode-api"]

FROM debian:bookworm-slim AS worker
RUN apt-get update \
    && apt-get install --yes --no-install-recommends ca-certificates ffmpeg \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --create-home recode \
    && mkdir -p /data \
    && chown recode:recode /data
COPY --from=build /out/worker /usr/local/bin/recode-worker
USER recode
ENTRYPOINT ["/usr/local/bin/recode-worker"]
