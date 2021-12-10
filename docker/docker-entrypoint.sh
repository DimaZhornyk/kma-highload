#!/usr/bin/env sh
set -ue

set -x

case $1 in
  test)
    exec go test ./...
    ;;
esac

exec "$@"