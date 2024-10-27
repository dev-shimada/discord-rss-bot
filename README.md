[![Go Report Card](https://goreportcard.com/badge/github.com/dev-shimada/discord-rss-bot)](https://goreportcard.com/report/github.com/dev-shimada/discord-rss-bot)
[![CI](https://github.com/dev-shimada/discord-rss-bot/actions/workflows/ci.yaml/badge.svg)](https://github.com/dev-shimada/discord-rss-bot/actions/workflows/ci.yaml)
[![build](https://github.com/dev-shimada/discord-rss-bot/actions/workflows/build.yaml/badge.svg)](https://github.com/dev-shimada/discord-rss-bot/actions/workflows/build.yaml)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)

# discord-rss-bot

## Getting started
```
docker compose up -d
```

## Usage
- `/subscribe <URL>`
- `/list`
- `/unsubscribe <ID>`

## Docker build
```console
docker build . -t discord-rss-bot
```

## Docker run
```console
docker run --rm --name discord-rss-bot -e DISCORD_BOT_TOKEN="" --mount type=bind,source="$(pwd)"/sqlite,target=/app/sqlite discord-rss-bot
```
