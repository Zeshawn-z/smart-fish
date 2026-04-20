#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

if [ ! -d .venv ]; then
  python3 -m venv .venv
fi

source .venv/bin/activate
pip install -r requirements.txt

HOST="${HOST:-0.0.0.0}"
PORT="${PORT:-8001}"

uvicorn app.main:app --host "$HOST" --port "$PORT"
