# Modtorio
Search mods.factorio.com from the command line! Intended for server environments where your "friends" constantly pester you to install new mods to your already poor, mod-overloaded server.

## Development requirements
* Go v1.16+

## Dependencies
* golang.org/x/crypto/ssh/terminal (used to hide user passwords)

## Compiling
1. Clone the project.
2. Install dependencies with `go get`.
3. `go build` or (`go install`).

## Running
`./modtorio help` (or `modtorio help` if you installed it).

## Licence
BSD-3 clause
