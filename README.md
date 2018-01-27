# SMBF Discord Bot

[![Maintainability](https://api.codeclimate.com/v1/badges/338e52b1d57e4df881aa/maintainability)](https://codeclimate.com/github/heshoots/discordbot/maintainability) [![Build Status](https://travis-ci.org/heshoots/discordbot.svg?branch=master)](https://travis-ci.org/heshoots/discordbot)

## Getting Started

### Building
```bash
# get dependencies
dep ensure
# go build -o discordbot
go build
```

### Running
```bash
./discordbot
# or
go run main.go
```

You will need a bunch of environment variables set in order to run this properly. The bot will not run without the DISCORD_BOT_DISCORD_API environment variable.
