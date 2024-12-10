<h1 align="center">
<img src="./apps/frontend/public/logo.svg" /><br />
Odin
</h1>

<h4 align="center">Enjoy movies and TV shows.</h4>

<h2 align="center">To be used with <a href="https://github.com/ad-on-is/odin-tv">Odin TV</a></h2>

![screenshot](./screenshots/odin-screenshot.png)

# Introduction

Odin is a self-hostable solution to watch movies and TV shows.

[] Test
[x] Best

## Key features

- Discover movies and shows
- Scrobble
- See your watchlists
- Setup custom sections from Trakt lists
- Scrape Jackett for Torrents
- Unrestrict links with RealDebrid

## Prerequisites

- A Trakt API account
  - Create a new App: <https://trakt.tv/oauth/applications/new>
  - Not down the Trakt `clientId` and `clientSecret`
- TMDB account
  - Note down the `apiKey`
- RealDebrid Account

## Setup with Docker (docker-compose)

```yml
version: "3"
services:
  odin-scraper:
    image: ghcr.io/ad-on-is/odin-movieshow/scraper
    container_name: odin-scraper
    restart: unless-stopped
    environment:
      - JACKETT_URL=http://jackett:9117
      - JACKETT_KEY=xxxxx

  odin-backend:
    image: ghcr.io/ad-on-is/odin-movieshow/backend
    container_name: odin-backend
    restart: always
    environment:
      - LOG_LEVEL=debug
    volumes:
      - ./pb_data:/pb_data

  odin-frontend:
    image: ghcr.io/ad-on-is/odin-movieshow/frontend
    container_name: odin-frontend
    restart: always
    environment:
      - NUXT_PB_URL=https://odin-backend.example.com # URL must be accessible within your network

  jackett:
    image: lscr.io/linuxserver/jackett:latest
    container_name: jackett
    environment:
      - TZ=Etc/UTC
      - AUTO_UPDATE=true
    volumes:
      - ./jackett:/config
    restart: unless-stopped
```

## Configuration

- Login as Admin
  - **User:** <admin@odin.local>, **Password:** adminOdin1
  - Connect to RealDebrid
  - Fill in the inputs for Trakt, TMDB, etc.
  - Create a user
- Login as user and connect to your Trakt user account

## Connecting to Odin TV

- Open Odin TV on your Android TV box
- If not already, login as your user in the Odin frontend, and go to devices
- Click on **Link device** and enter the code shown on your TV

## Running local dev environment

```bash
# install Bun
curl -fsSL https://bun.sh/install | bash

# lone the repo
git clone https://github.com/ad-on-is/odin-movieshow
cd odin-movieshow

# install dependencies
bun install

# create apps/backend/.env and provide JACKETT_URL and JACKETT_KEY

# run dev
bun run dev
```

## License

MIT

---

> [adisdurakovic.com](https://adisdurakovic.com) &nbsp;&middot;&nbsp;
> GitHub [@ad-on-is](https://github.com/ad-on-is) &nbsp;&middot;&nbsp;
> Built using [pocketbase](https://pocketbase.io/) and [Nuxt](https://nuxt.com/)
