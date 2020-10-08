# [JamDrop](https://jamdrop.app)

[JamDrop](https://jamdrop.app) is a webapp for queue sharing. Drop songs into your friends' queues! RickRoll them until y'all're nemeses! It's up to you!

## Contributing

### Overview

**Backend**
- Go
- Google Cloud Run (serverless Docker image executor) for hosting
- Google Cloud Build for Docker image build / push
- Google Firestore for persistence
- Google Cloud Scheduler for regular jobs
- Cloudflare for DNS

**Frontend**
- HTML/JS/SCSS
- Mithril for DOM rendering'n'stuf
- Parcel for bundling

### Getting Started

#### A few heads ups

- Currently using cookies as an auth method. That is not the right way to do auth lol.
- Even when running locally, you need to connect to Google Firestore to get to the development collections. [It appears there are some solutions for this](https://hub.docker.com/r/mtlynch/firestore-emulator/) which should be looked into.
- There's a bunch of long-polling going on from the frontend. Unfortunately, at this point Cloud Run does not provide websocket support, so long-polling is the best option for "live updating" stuffz. I do think, however, there should probably only be one endpoint being polled for updates, instead of the two that currently exist.

#### A to-do list

- [ ] Clone it
- [ ] Get a `gcloud-key.json` and put it at the root of your repo
- [ ] Get a `.env` file and put it at the root of your repo
- [ ] Get the backend working
  - [ ] Get the right version of Go working (I use [`gvm`](https://github.com/moovweb/gvm))
  - [ ] Build it: `$ go build cmd/server`
  - [ ] I like to use the [`gin`](https://github.com/codegangsta/gin) hot reloader (along with `scripts/dev/gin.sh`)
- [ ] Get the frontend working
  - [ ] `$ yarn`
  - [ ] `$ yarn watch`
  
Once you can get two build servers running (gin and parcel) you should be good to hit `http://localhost:8080`.

## Attributions

Icons made by [Freepik](https://www.flaticon.com/authors/freepik) from www.flaticon.com:
- [Jam icon](https://jamdrop.app/public/jam.489cda7e.svg)
