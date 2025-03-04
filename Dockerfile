FROM node:22.14.0 as bundle-builder
COPY public /public
WORKDIR /public
RUN npm ci && npm run bundle

FROM golang:1.24.1 as go-builder
ARG MTUI_VERSION="docker-dev"
WORKDIR /data
COPY go.* /data/
RUN go mod download
COPY . /data
COPY --from=bundle-builder /public/js/bundle* /data/public/js/
COPY --from=bundle-builder /public/node_modules /data/public/node_modules
RUN CGO_ENABLED=1 go build -ldflags="-s -w -extldflags=-static -X mtui/app.Version=$MTUI_VERSION" .

FROM alpine:3.21.3
COPY --from=go-builder /data/mtui /bin/mtui
RUN apk update && apk add git
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]