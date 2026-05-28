FROM node:22.22.3 as bundle-builder
COPY frontend /frontend
WORKDIR /frontend
RUN npm ci && npm run build

FROM golang:1.26.2 as go-builder
ARG MTUI_VERSION="docker-dev"
WORKDIR /data
COPY go.* /data/
RUN go mod download
COPY . /data
COPY --from=bundle-builder /frontend/dist/* /data/frontend/dist/
RUN CGO_ENABLED=1 go build -ldflags="-s -w -extldflags=-static -X mtui/app.Version=$MTUI_VERSION" .

FROM alpine:3.23.4
COPY --from=go-builder /data/mtui /bin/mtui
RUN apk update && apk add git
EXPOSE 8080
ENTRYPOINT ["/bin/mtui"]