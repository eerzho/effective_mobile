`CONFIG_PATH=./config/local.yaml` - path to config

`go run ./cmd/effective_mobile` - start app

`go run ./cmd/migrator` - migrations up 

`go run ./cmd/migrator --down` - migrations down

`swag init -d cmd/effective_mobile,internal/http,internal/lib,internal/domain` - generate docs

[Swagger](http://localhost:8000/swagger/index.html)