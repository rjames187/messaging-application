#!/usr/bin/env bash

export HOST_NAME=https://localhost:443

envsubst < config.js.template > config.js