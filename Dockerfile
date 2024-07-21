FROM node:20.15.0 as bundle-builder
COPY public /public
WORKDIR /public
RUN npm ci && npm run bundle

FROM golang:1.22.4 as go-builder
ARG MTUI_VERSION="docker-dev"
COPY . /data
COPY --from=bundle-builder /public/js/bundle* /data/public/js/
COPY --from=bundle-builder /public/node_modules /data/public/node_modules
WORKDIR /data
RUN CGO_ENABLED=1 go build -ldflags="-s -w -extldflags=-static -X mtui/app.Version=$MTUI_VERSION" .

FROM alpine:3.20.1
COPY --from=go-builder /data/mtui /bin/mtui
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]