Development

Prerequisites:
* docker
* docker-compose

Starting:
```sh
# builds the webapp and starts the vite dev server on http://localhost:5173
docker-compose up ui_webapp
# start the backend (embeds the built webapp from `public/dist`)
docker-compose up ui
```

The vite dev server (hot-reload) proxies `/api` (including the websocket proxy for the wasm client)
to the backend, set `API_TARGET` to change the backend url (default: `http://localhost:8080`).

Webapp only (without docker):
```sh
cd public
npm ci
npm run dev     # dev server
npm run build   # production build to `public/dist`, embedded into the go binary
npm run eslint
```
