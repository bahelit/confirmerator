#!/usr/bin/env bash

## Create user no home no shell, can't login.
adduser --system --group --no-create-home --shell /bin/false nats
adduser --system --group --no-create-home --shell /bin/false mongodb
adduser --system --group --no-create-home --shell /bin/false confirmerator
