#!/usr/bin/env bash

cd "$1" || exit

# Make go fetch directly from github for latest commit.
export GOPROXY=direct

go get -u github.com/flixurapp/flixur/pluginkit@main
go get -u github.com/flixurapp/flixur/proto/go@main
