# Go app template

<img alt="go-chi icon" src="https://i.ibb.co/SPZTkTX/chi.png" height="100px" />
<img alt="sqlite icon" src="https://i.ibb.co/SVPQmm5/sqlite.png" height="100px" />
<img alt="Nix icon" src="https://i.ibb.co/Xj6HfmC/nix.png" height="100px" />

Template repository for [nix](https://nixos.org/) flake + [go](https://go.dev) + [chi](https://go-chi.io/) + [sqlc](https://sqlc.dev/) + [go-migrate](https://github.com/golang-migrate/migrate).

## Develop

```sh
# Enter development environment using nix flake:
nix develop

# Create application configuration using example .env file
cp .env.dev.example .env.dev
ln -sv .env.dev .env

# Start application with live-reload
echo *.go | entr -r go run main.go
```

## Database

Load .env file:

```sh
. .env
```

Migrate database (using CLI):

```sh
migrate -path ./migrations -database "sqlite3://$DB_PATH" up
```

In case of `Dirty database`-errors:

```sh
migrate -path ./migrations -database "sqlite3://$DB_PATH" up
# error: Dirty database version 2. Fix and force version.

migrate -path ./migrations -database "sqlite3://$DB_PATH" force 1
migrate -path ./migrations -database "sqlite3://$DB_PATH" up
```

Generate go-bindings:

```sh
sqlc generate
```

## License

<p xmlns:cc="http://creativecommons.org/ns#" xmlns:dct="http://purl.org/dc/terms/"><a property="dct:title" rel="cc:attributionURL" href="https://github.com/dko1905/go-app-template">go-app-template</a> by <a rel="cc:attributionURL dct:creator" property="cc:attributionName" href="https://0chaos.eu">Daniel Florescu</a> is marked with <a href="https://creativecommons.org/publicdomain/zero/1.0/?ref=chooser-v1" target="_blank" rel="license noopener noreferrer" style="display:inline-block;">CC0 1.0<img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1" alt=""><img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/zero.svg?ref=chooser-v1" alt=""></a></p> 
