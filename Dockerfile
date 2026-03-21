# ---- Web client build ----
FROM oven/bun:1-alpine AS webclient

WORKDIR /web

RUN --mount=type=bind,target=/web,rw \
    --mount=type=cache,target=/web/pkg/web/@vervstack/matreshka/node_modules \
    --mount=type=cache,target=/web/pkg/web/Matreshka-UI/node_modules \
    --mount=type=cache,target=/root/.bun \
    cd /web/pkg/web/@vervstack/matreshka && \
    bun install --frozen-lockfile && \
    bun run build && \
    cd /web/pkg/web/Matreshka-UI && \
    bun install --frozen-lockfile && \
    bun run build && \
    mv dist /dist

# ---- Go app build ----
FROM --platform=$BUILDPLATFORM golang:1.24.2-alpine AS builder

WORKDIR /src

COPY --from=webclient /dist /dist

RUN --mount=type=bind,target=/src,rw \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/go/mod \
    mkdir -p /src/internal/transport/web/dist && \
    mv /dist/* /src/internal/transport/web/dist && \
    go mod download && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 \
    go build -o /deploy/server/service /src/cmd/service/main.go && \
    cp -r config /deploy/server/config && \
    mkdir -p /deploy/server/migrations && \
    cp -r /src/migrations/* /deploy/server/migrations/

FROM alpine:3.14

LABEL MATRESHKA_CONFIG_ENABLED=true

WORKDIR /app

COPY --from=builder /deploy/server/ .

RUN mkdir /app/data

EXPOSE 50049

VOLUME './data/matreshka-be.db'

ENTRYPOINT ["./service"]