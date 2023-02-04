#!/bin/bash
for service in trainer trainings users; do
  /usr/local/bin/docker-entrypoint.sh generate \
    -i "./openapi/$service.yml" \
    -g javascript \
    -o "./web/src/repositories/clients/$service"
done
