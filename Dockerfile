FROM node:20.18.0 as bundle-builder
COPY public /public
WORKDIR /public
RUN npm ci && npm run bundle

FROM golang:1.23.2 as go-builder
ARG MTUI_VERSION="docker-dev"
WORKDIR /data
COPY go.* /data/
RUN go mod download
COPY . /data
COPY --from=bundle-builder /public/js/bundle* /data/public/js/
COPY --from=bundle-builder /public/node_modules /data/public/node_modules
RUN CGO_ENABLED=1 go build -ldflags="-s -w -extldflags=-static -X mtui/app.Version=$MTUI_VERSION" .

FROM alpine:3.20.3
COPY --from=go-builder /data/mtui /bin/mtui
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]