# mtui

Minetest web ui

![](https://github.com/minetest-go/mtui/workflows/test/badge.svg)
![](https://github.com/minetest-go/mtui/workflows/build/badge.svg)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minetest-go/mtui)
[![Go Report Card](https://goreportcard.com/badge/github.com/minetest-go/mtui)](https://goreportcard.com/report/github.com/minetest-go/mtui)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/mtui/badge.svg)](https://coveralls.io/github/minetest-go/mtui)

# Features

* Account/Password management
* Remote console (chat, lua-injection)
* Chat console
* World status
* Skin management
* Mail management
* XBan management
* Mediaserver (remote_media)
* OAuth provider
* Event logging
* Mod/game/texturepack configuration and updates (cdb, git)
* Engine management (via Docker)
* File browser
* Minetest config management

Planned:
* Mapserver management (via Docker)
* Matterbridge management (via Docker)

# Running

## Environment Variables

* `WORLD_DIR` world directory, defaults to the current working dir
* `WEBDEV` if set to "true": bypasses the embedded web-resources (for development)
* `API_KEY` api key, optional, will be generated if not set
* `COOKIE_DOMAIN` the cookie domain, defaults to "127.0.0.1"
* `COOKIE_PATH` the cookie path, defaults to "/"
* `COOKIE_SECURE` secure cookie, defaults to "false"
* `ADMIN_USERNAME` initial admin username (optional)
* `ADMIN_PASSWORD` initial admin password (optional)
* `LOGLEVEL` currently supported: "debug", default is "info"
* `SERVER_NAME` Server-name to display on the ui
* `ENABLE_FEATURES` manually enabled features
* `MINETEST_CONFIG` set this to the `minetest.conf` location to enable the settings-management

# Api

* See [http-api](./http-api.md)

# Development

Prerequisites:
* docker
* docker-compose

Starting:
```sh
# init and update submodules
git submodule init
git submodule update
# fetch the frontend libraries (one time task)
docker-compose up ui_webapp
# start the minetest engine and the ui
docker-compose up ui minetest
```

# License

* Code: `MIT`
* Textures: `CC BY-SA 3.0`
  * `public/pics/sam.png` [minetest_game](https://github.com/minetest/minetest_game)