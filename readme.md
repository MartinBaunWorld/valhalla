# Valhalla - Inheritance

When you die, you might want to leave a message or a secret for your loved ones.
This is the simple way to do it.

## How it works

1. Setup a telegram bot and group - your
2. Setup environment variables - see example below
3. Run it

That's it.

Checkout how I use this app to ensure my [employees inherit some Bitcoin](https://..) 

## Example Run

You can make a script that runs


```
#!/bin/bash

set -e

export TELEGRAM_KEY="mock"
export TELEGRAM_CHAT="TODO"
export WEB_URL='http://localhost:8000'
export SECRET="Valhalla is here"
export SEND_SECRET_AFTER=3
export ASK_EVERY_SECONDS=3
# export ASK_EVERY_DAYS=3

./valhalla.py
```

Replace the environment variables. The `ASK_EVERY_SECONDS=3` can be helpful for testing, and then later you can set `ASK_EVERY_DAY=3`
