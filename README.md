# discord-rss-bot

<p align="center">
  <a href="https://github.com/dev-shimada/discord-rss-bot/actions/workflows/ci.yaml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/dev-shimada/discord-rss-bot/ci.yaml?branch=main&style=flat-square&logo=githubactions&logoColor=white&label=CI"></a>
  <a href="https://coveralls.io/github/dev-shimada/discord-rss-bot?branch=main"><img alt="Coverage" src="https://img.shields.io/coverallsCoverage/github/dev-shimada/discord-rss-bot?branch=main&style=flat-square&logo=coveralls&logoColor=white"></a>
  <a href="https://github.com/dev-shimada/discord-rss-bot/blob/main/go.mod"><img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/dev-shimada/discord-rss-bot?style=flat-square&logo=go&logoColor=white"></a>
  <a href="https://github.com/dev-shimada/discord-rss-bot/blob/main/LICENSE"><img alt="MIT License" src="https://img.shields.io/github/license/dev-shimada/discord-rss-bot?style=flat-square"></a>
</p>

<p align="center">
  <a href="https://github.com/dev-shimada/discord-rss-bot/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/dev-shimada/discord-rss-bot?style=flat-square&logo=github&logoColor=white"></a>
  <a href="https://github.com/dev-shimada/discord-rss-bot/actions/workflows/build.yaml"><img alt="Docker Build" src="https://img.shields.io/github/actions/workflow/status/dev-shimada/discord-rss-bot/build.yaml?style=flat-square&logo=docker&logoColor=white&label=build"></a>
  <a href="https://goreportcard.com/report/github.com/dev-shimada/discord-rss-bot"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/dev-shimada/discord-rss-bot?style=flat-square"></a>
</p>

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
