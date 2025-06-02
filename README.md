# Zulip to Telegram gate

This bot forwards messages from Zulip, in which keywords are mentioned, to a Telegram Bot.

# Setup

The configuration setup is intuitive and requires no special explanations.

```
tg-key: 1234567890:AQGn9v_SmFJ05w3PF0yhBFi6lVSiOFgmqjQ
bot-user-id: 1234567890

zulip:
  site: https://zulip.server.by

users:
  - name: ivanov
    tg-id: 234567891
    zulip-email: ivanov@example.com
    zulip-key: OCdORmfVbfCltmitWD2j69LXTEfiyBQr
    words:
      - Ivanov
  - name: petrov
    tg-id: 345678901
    zulip-email: petrov@example.com
    zulip-key: N4Zp9d5NBJqytlm0bP1JVwwcjA1X9APk
    words:
      - Petrov
      - Petroff
    ignore-private-from:
      - angry-bot@example.com
      - nospam@example.com
```

Use BotFather in Telegram to create a bot in which you will receive messages from the Zulip messenger. Open a chat with the created bot and run the command /start.

In the settings of your profile in the Zulip messenger, open the Account & privacy section. In the API key section, click Manage your API key and create a token.

Then compile and run.
