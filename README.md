<h1 align="center">
<img src="./apps/frontend/public/logo.svg" />
</h1>

<h4 align="center">Enjoy movies and TV shows.</h4>

<h2 align="center">To be used with <a href="https://github.com/ad-on-is/odin-tv">Odin TV</a></h2>

![screenshot](./screenshots/odin-screenshot.png)

# Introduction

Odin is a self-hostable solution to watch movies and TV shows.

## Key features

- Discover movies and shows
- Scrobble
- See your watchlists
- Setup custom sections from Trakt lists
- Scrape Jackett for Torrents
- Unrestrict links with RealDebrid/AllDebrid

## Prerequisites

- A Trakt API account
  - Create a new App: <https://trakt.tv/oauth/applications/new>
  - Note down the Trakt `clientId` and `clientSecret`
- TMDB account
  - Note down the `apiKey`
- One of:
  - RealDebrid Account
  - AllDebrid Account

## Setup with Docker (docker-compose)

```yaml
services:
  odin-backend:
    image: ghcr.io/ad-on-is/odin-movieshow/backend
    container_name: odin-backend
    restart: always
    environment:
      - LOG_LEVEL=debug
      - JACKETT_URL=http://jackett:9117
      - JACKETT_KEY=xxxxx
      - BACKEND_URL=http://odin-backend.example.com # URL must be accessible within your network
      - MQTT_URL=wss://mqtt.example.com # URL must be accessible within your network
      - MQTT_USER=<user>
      - MQTT_PASSWORD=<password>
      - TMDB_KEY=<tmdbkey>
      - TRAKT_CLIENTID=<trakt_clientid>
      - TRAKT_SECRET=<trakt_secret>
      - ALLDEBRID_KEY=<alldebrid_key>
    volumes:
      - ./pb_data:/pb_data

  odin-frontend:
    image: ghcr.io/ad-on-is/odin-movieshow/frontend
    container_name: odin-frontend
    restart: always
    environment:
      - BACKEND_URL=https://odin-backend.example.com # URL must be accessible within your network
      - MQTT_URL=wss://mqtt.example.com # URL must be accessible within your network
      - MQTT_USER=<user>
      - MQTT_PASSWORD=<password>

  # this is just an example config for Jackett
  jackett:
    image: lscr.io/linuxserver/jackett:latest
    container_name: jackett
    environment:
      - TZ=Etc/UTC
      - AUTO_UPDATE=true
    volumes:
      - ./jackett:/config
    restart: always
```

## Configuration

- Login as Admin

  - **User:** <admin@odin.local>, **Password:** adminOdin1
  - Connect to RealDebrid by following the steps

  ## Creating a user

- Create a new user
- Login as user and connect to your Trakt user account, by following the steps

## Connecting to Odin TV

> [!NOTE]
> This only works with a regular user, not an admin account.

- Open Odin TV on your Android TV box
- If not already, login as your user in the Odin frontend, and go to devices
- Click on **Link device**, enter the code shown on your TV and click **Connect**

## Running local dev environment

```bash
# install Bun
curl -fsSL https://bun.sh/install | bash

# lone the repo
git clone https://github.com/ad-on-is/odin-movieshow
cd odin-movieshow

# install dependencies
bun install

# copy .env.example to apps/backend/.env and apps/frontend/.env and fill in the blanks

# run dev
bun --bun run dev
```

## License

MIT

---

> GitHub [@ad-on-is](https://github.com/ad-on-is) &nbsp;&middot;&nbsp;
> Built using [pocketbase](https://pocketbase.io/) and [Nuxt](https://nuxt.com/)
