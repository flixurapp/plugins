#!/usr/bin/env bash

cd "$1" || exit

# Make go fetch directly from github for latest commit.
export GOPROXY=direct

go get -u forge.xela.codes/xela/flixur/pluginkit@master
go get -u forge.xela.codes/xela/flixur/proto/go@master
