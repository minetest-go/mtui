 Development

Prerequisites:
* docker
* docker-compose

Starting:
```sh
# download npm packages
docker compose up npm
# process 1: frontend build/watch
docker compose up watch
# process 2: backend api
docker compose up ui
```