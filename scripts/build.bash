#! /usr/bin/env bash

commit_hash=$(git rev-list -1 HEAD)
commit_message=$(git log -1 --pretty=%B | cat)
build_time=$(date +'%FT%T')
linker_flags="\
    -X 'jamdrop/config.CommitHash=$commit_hash'
    -X 'jamdrop/config.CommitMessage=$commit_message' 
    -X 'jamdrop/config.BuildTime=$build_time'"

go build -ldflags "$linker_flags" cmd/server/main.go
