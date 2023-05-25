
Accessing the mtui with `curl`

# Authentication

## Login

Login with a cookie-jar to store the session cookie:

```sh
curl --cookie cookies --cookie-jar cookies --data '{"username":"user","password":"pass"}' -H "Content-Type: json" http://127.0.0.1:8080/api/login
```

With otp:

```sh
curl --cookie cookies --cookie-jar cookies --data '{"username":"user","password":"pass","otp_code":"123456"}' -H "Content-Type: json" http://127.0.0.1:8080/api/login
```


Response:

Success (status 200):
```json
{
  "exp": 1685622561,
  "username": "user",
  "privileges": [
    "interact",
    "server"
  ]
}
```

Error-status:
* `401` invalid password
* `403` invalid otp
* `404` user not found

## Refresh / Check login

Refresh or check the login:

```sh
curl --cookie cookies --cookie-jar cookies http://127.0.0.1:8080/api/login
```

Response:

```json
{
  "exp": 1685622561,
  "username": "user",
  "privileges": [
    "interact",
    "server"
  ]
}
```


# Execute commands

**NOTE:** Executing commands requires an authenticated session (see: login)

## Chat command

Execute `/status`:

```sh
curl -i --cookie cookies --cookie-jar cookies --data '{"playername":"user","command":"status"}' -H "Content-Type: application/json" http://127.0.0.1:8080/api/bridge/execute_chatcommand
```

Response:

```json
{
  "success": true,
  "message": "# Server: version: 5.7.0 | game: Minetest Game | uptime: 12min 27s | max lag: 0.323s | clients: "
}
```

Error-response (also with status 200):

```json
{
  "success": false,
  "message": "invalid command"
}
```

## Lua command

Return `minetest.features` table:

```sh
curl --cookie cookies --cookie-jar cookies --data '{"code": "return minetest.features"}' -H "Content-Type: application/json" http://127.0.0.1:8080/api/bridge/lua
```

**NOTE**: the return value is returned in json format in the `result` field

Success-response:

```json
{
  "success": true,
  "message": "",
  "result": {
    "abm_min_max_y": true,
    "add_entity_with_staticdata": true,
    "area_store_custom_ids": true,
    "area_store_persistent_ids": true,
    "compress_zstd": true,
    "degrotate_240_steps": true,
    "direct_velocity_on_players": true,
    "dynamic_add_media_table": true,
    "formspec_version_element": true,
    "get_all_craft_recipes_works": true,
    "get_light_data_buffer": true,
    "get_sky_as_table": true,
    "glasslike_framed": true,
    "httpfetch_binary_data": true,
    "mod_storage_on_disk": true,
    "no_chat_message_prediction": true,
    "no_legacy_abms": true,
    "nodebox_as_selectionbox": true,
    "object_independent_selectionbox": true,
    "object_step_has_moveresult": true,
    "object_use_texture_alpha": true,
    "particlespawner_tweenable": true,
    "pathfinder_works": true,
    "texture_names_parens": true,
    "use_texture_alpha": true,
    "use_texture_alpha_string_modes": true
  }
}
```

Error (status 200):

```json
{
  "success": false,
  "message": "Command crashed: \"...inetest/worlds/world/worldmods/mtui_mod/handlers/lua.lua:7: attempt to call upvalue 'fn' (a nil value)\"",
  "result": null
}
```

# Player infos

Request:

```sh
curl --cookie cookies --cookie-jar cookies http://127.0.0.1:8080/api/player/info/user
```

Response (status 200):

```json
{
  "auth_entry": true,
  "auth_id": 1102,
  "player_entry": true,
  "name": "user",
  "privs": [
    "fast",
    "home",
    "interact",
    "shout",
    "tp",
    "zoom"
  ],
  "last_login": 1685017281,
  "first_login": 1636904045,
  "breath": 10,
  "health": 19,
  "pitch": 45.0999985,
  "yaw": 214.800003,
  "posx": -4353.72021,
  "posy": 230.630005,
  "posz": -761.909973,
  "stats": {
    "crafted": 277,
    "died": 22,
    "digged_nodes": 1948,
    "placed_nodes": 1660,
    "played_time": 122992
  }
}
```

# XBan records

Request:

```sh
curl --cookie cookies --cookie-jar cookies http://127.0.0.1:8080/api/xban/records/user
```

Response::

```json
{
  "banned": false,
  "last_seen": 1685012244,
  "names": {
    "127.0.0.1": true,
    "user": true
  },
  "record": null
}
```