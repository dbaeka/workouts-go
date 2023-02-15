include .env

.PHONY: lint
.PHONY: proto
proto:
	protoc --go_out=internal/common/genproto/trainer --go_opt=paths=source_relative \
--go-grpc_out=internal/common/genproto/trainer --go-grpc_opt=paths=source_relative -I api/protobuf api/protobuf/trainer.proto
	protoc --go_out=internal/common/genproto/users --go_opt=paths=source_relative \
  --go-grpc_out=internal/common/genproto/users --go-grpc_opt=paths=source_relative -I api/protobuf api/protobuf/users.proto

lint:
	@./scripts/lint.sh trainer
	@./scripts/lint.sh trainings
	@./scripts/lint.sh users

.PHONY: openapi
openapi: openapi_http openapi_js

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o internal/trainings/ports/openapi_types.gen.go -package ports api/openapi/trainings.yml
	oapi-codegen -generate chi-server -o internal/trainings/ports/openapi_api.gen.go -package ports api/openapi/trainings.yml

	oapi-codegen -generate types -o internal/trainer/ports/openapi_types.gen.go -package ports api/openapi/trainer.yml
	oapi-codegen -generate chi-server -o internal/trainer/ports/openapi_api.gen.go -package ports api/openapi/trainer.yml

	oapi-codegen -generate types -o internal/users/openapi_types.gen.go -package main api/openapi/users.yml
	oapi-codegen -generate chi-server -o internal/users/openapi_api.gen.go -package main api/openapi/users.yml

.PHONY: openapi_js
openapi_js:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
  -i /local/api/openapi/trainings.yml \
  -g javascript \
  -o /local/web/src/repositories/clients/trainings

	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
  -i /local/api/openapi/trainer.yml \
  -g javascript \
  -o /local/web/src/repositories/clients/trainer

	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
  -i /local/api/openapi/users.yml \
  -g javascript \
  -o /local/web/src/repositories/clients/users

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