#!/usr/bin/env python3
import os
import threading
import requests
from time import sleep
from traceback import format_exc
from datetime import timedelta, datetime 

from bottle import get, post, run


def get_str(name):
    if name in os.environ:
        return os.environ[name]

    raise Exception(f'Expected env variable {name} to be set')


def get_int(name, default=None):
    try:
        return int(get_str(name))
    except:
        if default is not None:
            return int(default)
        raise Exception(f"Expected {name} to be int")

TELEGRAM_KEY = get_str('TELEGRAM_KEY')
TELEGRAM_CHAT = get_str('TELEGRAM_CHAT')
WEB_URL = get_str('WEB_URL')
SECRET = get_str('SECRET')

SEND_SECRET_AFTER = get_int('SEND_SECRET_AFTER')
ASK_EVERY = timedelta(days=get_int('ASK_EVERY_DAYS', 0), seconds=get_int('ASK_EVERY_SECONDS', 0))


if int(ASK_EVERY.total_seconds()) == 0:
    print('You need to define ASK_EVERY_DAYS or ASK_EVERY_SECONDS')
    exit(1)


asked = 0
reset_at = datetime.now()

def msg(msg):
    print(msg)
    if TELEGRAM_KEY == "mock":
        return
    try:
        res = requests.post(
            f"https://api.telegram.org/bot{TELEGRAM_KEY}/sendMessage",
            json=dict(chat_id=TELEGRAM_CHAT, text=msg),
        )
        if res.status_code != 200:
            print(res.text)
        return True

    except KeyboardInterrupt:
        print("exiting by keyboard interrupt...")
        exit(0)
    except:  # noqa
        print(format_exc())
        return False


def checker():
    global reset_at, asked
    while True:
        sleep(ASK_EVERY.total_seconds())
        if asked == SEND_SECRET_AFTER:
            msg(SECRET)
            exit(0)
        msg(f"Click here if Valhalla is awaiting\n{WEB_URL}/valhalla")
        asked += 1


@get("/valhalla")
def valhalla():
    with open('index.html', 'r') as f:
        return f.read()


@post("/valhalla/reset_timer")
def reset_timer():
    global reset_at, asked
    reset_at = datetime.now()
    asked = 0
    msg("Timer Reset")
    return "Timer reset"


msg("Started Valhalla")
x = threading.Thread(target=checker, daemon=True)
x.start()
run(port=8000, host='0.0.0.0')
msg("Finished Valhalla")
