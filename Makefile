# If the first argument is "migrate"...
ifeq (migrate, $(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "migrate"
  MIGRATE_ARGS := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(MIGRATE_ARGS):;@:)
endif

.PHONY: migrate
migrate:
	go run cmd/migrate/main.go \
	-source file://migration \
	-database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable \
	$(MIGRATE_ARGS)

# Docker commands
.PHONY: docker-build docker-build-no-cache docker-compose-up

# Сборка Docker образа
docker-build: vendor
	docker build --build-arg BUILDKIT_INLINE_CACHE=1 -t event-service:latest .

# Сборка без кэша
docker-build-no-cache:
	docker build --no-cache -t service:latest .

docker-compose-up:
	docker-compose up -d

.PHONY: swagger-gen
swagger-gen: swagger-fmt
	swag init \
        --parseDependency \
        --parseInternal \
        --parseDepth 5 \
        -g internal/app/start/start.go \
        --output docs/swagger

.PHONY: swagger-fmt
swagger-fmt:
	swag fmt

.PHONY: swagger-install
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3