help:
	@echo "\
Usage: \n\
    make \n\
    Commands: \n\
         | help                         Показать это сообщение\n\
         | build                        Собрать проект\n\
         | run                          Запустить проект\n\
         | stop                         Остановить проект\n\
         | restart                      Перезапустить проект\n\
         | lint                         Проверить код на ошибки\n\
         | lint-fix                     Исправить ошибки в коде\n\
         | test                         Запустить тесты\n\
         | migrate-create   	        Создать новую миграцию \n\
         | migrate-up                   Применить все миграции\n\
         | migrate-down                 Откатить последнюю миграцию\n\
         | migrate-status               Проверить статус миграций\n\
         | migrate-reset                Сбросить все миграции\n\
    "


migrate-create:
	@read -p "Введите имя миграции: " MIGRATION_NAME; \
	goose create "$$MIGRATION_NAME" sql

migrate-up:
	goose up

migrate-down:
	goose down

migrate-status:
	goose status

migrate-reset:
	goose reset

lint:
	golangci-lint run ./...

test:
	go test -race -count 10 ./internal/...

build:
	docker compose build

run:
	docker compose up -d

stop:
	docker compose down

restart:
	docker compose restart

lint-fix:
	gofmt -s -w .
	golangci-lint run --fix

