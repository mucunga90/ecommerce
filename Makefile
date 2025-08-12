.PHONY: test-db-up test-db-down test

POSTGRES_CONTAINER_NAME=test-postgres
POSTGRES_PORT=5433

test-db-up:
	docker run --rm -d \
        --name $(POSTGRES_CONTAINER_NAME) \
        -e POSTGRES_USER=test \
        -e POSTGRES_PASSWORD=test \
        -e POSTGRES_DB=testdb \
        -p $(POSTGRES_PORT):5432 \
        postgres:16-alpine
	@echo "Postgres test DB started on port $(POSTGRES_PORT)"

test-db-down:
	docker stop $(POSTGRES_CONTAINER_NAME)
	docker rm -f $(POSTGRES_CONTAINER_NAME) || true

test: test-db-up
	@echo "Waiting for DB to be ready..."
	sleep 3
	@echo "Executing golang test"
	go test ./internal/storage
	sleep 3
	$(MAKE) test-db-down