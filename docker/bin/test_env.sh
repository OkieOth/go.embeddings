#!/bin/bash
scriptPos=${0%/*}

set -e

COMPOSE_FILE=$scriptPos/../test_env.yaml

docker compose -p go-embeddings-test -f $COMPOSE_FILE up -V --build --always-recreate-deps --abort-on-container-exit --exit-code-from test_runner
