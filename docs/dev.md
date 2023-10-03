 Development

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