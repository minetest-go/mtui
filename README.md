# mtui

Minetest web ui

![](https://github.com/minetest-go/mtui/workflows/test/badge.svg)
![](https://github.com/minetest-go/mtui/workflows/build/badge.svg)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minetest-go/mtui)
[![Go Report Card](https://goreportcard.com/badge/github.com/minetest-go/mtui)](https://goreportcard.com/report/github.com/minetest-go/mtui)
[![Coverage Status](https://coveralls.io/repos/github/minetest-go/mtui/badge.svg)](https://coveralls.io/github/minetest-go/mtui)

# Features

* Account/Password management
* Remote console
* World status
* skin management

Planned:
* mod/game/texturepack configuration and updates

# Running

## Configuration

Application-config: `<world_dir>/mtui.json` (will be automatically created)

Example:
```json
{
    "jwt_key": "mysecretjwtkey",
    "api_key": "mykey",
    "cookie_domain": "127.0.0.1",
    "cookie_secure": false,
    "cookie_path": "/"
}
```

## Environment Variables

* `WORLD_DIR` world directory, defaults to the current working dir
* `WEBDEV` if set to "true": bypasses the embedded web-resources (for development)

# Development

Prerequisites:
* docker
* docker-compose

Starting:
```sh
# fetch the frontend libraries (one time task)
docker-compose up ui_webapp
# start the minetest engine and the ui
docker-compose up ui minetest
```

# License

MIT