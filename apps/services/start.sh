#!/bin/bash

BLUE="\e[34m"
GREEN="\e[32m"
NC="\e[0m"

docker compose down && docker compose up -d

echo -e "

=================================================
${GREEN}AiO Server running on:${NC} http://localhost:6060
=================================================

${BLUE} ➜ Frontend:${NC} http://localhost:6060/
${BLUE} ➜ Backend:${NC} http://localhost:6060/api/
${BLUE} ➜ PocketBase:${NC} http://localhost:6060/_/
${BLUE} ➜ MQTT:${NC} http://localhost:6060/socket/

"
