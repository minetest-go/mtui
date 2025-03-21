
# installation

Environment variables:
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

# Using docker-compose
You must use file `docker-compose.yml` from this entire repository directory, because it pulls in the app files from this repository. You cannot run it in a bare directory.

Container `ui_webapp` builds the app, and does not run it. So you only need to do this until it finishes building:
```
docker-compose up ui_webapp
```
When the app is completely built, then you can run only container *ui*.
```
docker-compose up -d ui
```
