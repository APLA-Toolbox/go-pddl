# go-pddl

A PDDL Parser written in Go

# Purpose

The go-al is to parse a PDDL file and convert it to a JSON/YAML standard format.

# Run

- Install the latest version of go

## Docker

- Install docker
- Install docker-compose
- `go get github.com/APLA-Toolbox/go-pddl/edit/main/README.md`
- `cd $GOPATH/src/github.com/APLA-Toolbox/go-pddl`
- `docker-compose -f docker/docker-compose.yml up`

## Scripting

- Create a .env file, copy the content of .env.example
- Fill DOMAIN and PROBLEM with domain.pddl and problem.pddl absolute paths
- Source .env
- Run go run main.go

# Contributions

See the open issues and feel free to contribute, help is WANTED.
