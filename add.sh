#!/bin/bash

url=$1
go run ./src/add/cmd/reading-list add log.json "$url"
git commit -am .
git push origin HEAD
