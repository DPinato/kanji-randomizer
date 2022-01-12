#!/bin/bash
terraform init

# prepare executable for AWS lambda
GOOS=linux go build kanji-randomizer.go

zip kanji-randomizer.zip kanji-randomizer
terraform apply

# clean up
rm kanji-randomizer kanji-randomizer.zip