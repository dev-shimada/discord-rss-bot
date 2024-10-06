# discord-rss-bot

## Discord command
- `!subscribe url`

## Docker build
```console
docker build . -t discord-rss-bot
```

## Docker run
```console
docker run --rm --name discord-rss-bot -e DISCORD_BOT_TOKEN="" --mount type=bind,source="$(pwd)"/sqlite,target=/app/sqlite discord-rss-bot
```
