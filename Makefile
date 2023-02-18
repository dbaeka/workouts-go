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


INTERNAL_PACKAGES := $(wildcard internal/*)

ifeq (test,$(firstword $(MAKECMDGOALS)))
  TEST_ARGS := $(subst $$,$$$$,$(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)))
  $(eval $(TEST_ARGS):;@:)
endif
.PHONY: test $(INTERNAL_PACKAGES)
test: $(INTERNAL_PACKAGES)
$(INTERNAL_PACKAGES):
	@(cd $@ && go test -count=1 -race ./... $(subst $$$$,$$,$(TEST_ARGS)))

.PHONY: fmt
fmt:
	goimports -l -w internal/