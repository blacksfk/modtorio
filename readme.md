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

## Docker
1. `docker build -t modtorio .`
2. `docker run -i -t modtorio <command>`

Obviously the above is only really useful for searching and viewing the help documentation. In order to download, update, enable and disable, or list mods a directory must be provided.

`docker run -i -t -v /absolute/path/to/mod/dir:/mods modtorio --dir /mods <command>`

Providing a populated `mod-list.json` but no mod files in the directory will result in modtorio downloading the latest version for each mod in `mod-list.json`.

## Licence
BSD-3 clause
