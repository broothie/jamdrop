
linker_flags="-X 'jamdrop/config.BuildTime=$(date +'%FT%T')'"

go build -ldflags "$linker_flags" cmd/server/main.go
