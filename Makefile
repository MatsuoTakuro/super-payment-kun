
init: build run

build:
	docker compose build

run:
	docker compose up -d

down:
	docker compose down  --remove-orphans

clean: down rm_volume

rm_volume:
	docker volume rm super-payment-kun_db_data

test:
	go test -v -race -shuffle=on -covermode=atomic ./...

log_app:
	docker compose logs -f app

db_in:
	docker compose exec db mysql -u taro -ppass super-payment-kun-db

migrate-dry:
	mysqldef -u taro -p pass -h 127.0.0.1 -P 3316 super-payment-kun-db --dry-run < ./_tools/mysql/init/00_schema.sql

migrate:
	mysqldef -u taro -p pass -h 127.0.0.1 -P 3316 super-payment-kun-db < ./_tools/mysql/init/00_schema.sql

generate:
	go generate ./...

watch_api_spec: # install redocly/cli command first
	npx @redocly/cli preview-docs  doc/api_spec.yaml --host "127.0.0.1" --port 65535
