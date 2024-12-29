#!/bin/bash

node /odin-frontend/server/index.mjs &
mosquitto -c /etc/mosquitto/mosquitto.conf &
caddy run --config /etc/caddy/Caddyfile &
/odin-server serve --http=0.0.0.0:8090
