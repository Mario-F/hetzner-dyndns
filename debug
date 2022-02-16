#!/bin/sh
DEBUG_PORT=${DEBUG_PORT:-8345}
dlv debug --listen=:${DEBUG_PORT} --headless --log --api-version=2 -- "$@"
