include .env

.PHONY: lint
.PHONY: proto
proto:
	@./scripts/proto.sh trainer
	@./scripts/proto.sh users

lint:
	@go-cleanarch
	@./scripts/lint.sh common
	@./scripts/lint.sh trainer
	@./scripts/lint.sh trainings
	@./scripts/lint.sh users

.PHONY: openapi
openapi: openapi_http openapi_js

.PHONY: openapi_http
openapi_http:
	@./scripts/openapi-http.sh trainer internal/trainer/ports ports
	@./scripts/openapi-http.sh trainings internal/trainings/ports ports
	@./scripts/openapi-http.sh users internal/users main

.PHONY: openapi_js
openapi_js:
	@./scripts/openapi-js.sh trainer
	@./scripts/openapi-js.sh trainings
	@./scripts/openapi-js.sh users

.PHONY: mysql
mysql:
	mysql -u ${MYSQL_USERNAME} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}

.PHONY: test $(INTERNAL_PACKAGES)
test:
	@./scripts/test.sh common .e2e.env
	@./scripts/test.sh trainer .test.env
	@./scripts/test.sh trainings .test.env
	@./scripts/test.sh users .test.env

.PHONY: fmt
fmt:
	goimports -l -w internal/