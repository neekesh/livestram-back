#!/bin/bash

ginWatch() {
  PORT=$2
  if [ -z "$PORT" ]; then
    PORT=3000;
  fi
  gin -a $1 -i -p $PORT run .
}

if which gin; then
  ginWatch $1 $2
else
  go install github.com/codegangsta/gin@latest
  ginWatch $1 $2
fi
