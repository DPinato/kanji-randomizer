#!/bin/bash

# build the app
go mod tidy
go build -o kanji-randomizer kanji-randomizer.go

# create the /usr/local/bin/twitchjoiner directory and move executable there
TJ_DIR="/usr/local/bin/kanji-randomizer"
if [ ! -d $TJ_DIR ]; then
  sudo mkdir -p $TJ_DIR;
fi

sudo chown pi $TJ_DIR
cp -t $TJ_DIR/ kanji-randomizer
rm kanji-randomizer