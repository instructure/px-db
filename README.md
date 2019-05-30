# Practice Database Sanitizer (PostgreSQL)

## Setup

Consider the following: `git config --global url."git@github.com:".insteadOf "https://github.com/"`

- `go get -u github.com/instructure/px-db`
- OR Download the binary: Place the binary in `/usr/local/bin/` or in your Shell `$PATH`
- OR Pull the project and compile

## Building Locally

```bash
go get github.com/instructure/px-cg-deploy
cd $GOPATH/src/github.com/instructure/px-cg-deploy
go build
```

## Usage

```
px-db is for sanitizing Practice PostgreSQL Tables

Usage:
  px-db [flags]
  px-db [command]

Available Commands:
  display     Display data from a PostgreSQL DB
  help        Help about any command
  plugin      Run Plugins that perform custom logic for PostgreSQL DB Table Sanitization
  sanitize    Sanitize a PostgreSQL DB

Flags:
      --db-endpoint string   PostgreSQL Hostname/Endpoint
      --db-name string       PostgreSQL Database Name
      --db-port int          PostgreSQL Bind Port (default 5432)
      --db-ssl-mode          PostgreSQL SSL Mode (default false)
      --db-user string       PostgreSQL User
  -h, --help                 help for px-db
      --version              version for px-db

Use "px-db [command] --help" for more information about a command.
```

See the `hack` directory for the Practice implementation and application.
