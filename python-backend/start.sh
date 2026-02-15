#!/bin/sh

set -eu

ENV_NAME="${1:-local}"

if [ "$ENV_NAME" = "prod" ]; then
  python run.py --env prod --host 0.0.0.0 --port 8123
else
  python run.py --env "$ENV_NAME" --host 0.0.0.0 --port 8123 --reload
fi
