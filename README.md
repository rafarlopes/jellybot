# JellyBot

This is a prototype of a bot to run Docker commands using a chat bot.

### Env Variables

```sh

export CHAT_API_TOKEN=<YOUR_SLACK_API_TOKEN>
export BOT_USER_ID=<BOT_USER_ID_CREATED_ON_SLACK>
export DEFAULT_CHAT_ID=<CHANNEL_ID_THE_SHOULD_ANSWER>

```

### Deps

This project uses [dep](https://github.com/golang/dep) as dependency management tool.

To build, use as follow:

```sh

dep ensure
go build

```