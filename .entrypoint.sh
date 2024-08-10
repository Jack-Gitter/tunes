#!/bin/sh

go install github.com/pressly/goose/v3/cmd/goose@latest 
make goose-up 
./tunes

