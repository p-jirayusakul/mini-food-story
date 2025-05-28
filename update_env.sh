#!/bin/bash

ENV_FILE=".env"

KEY=$1
VALUE=$2

if grep -q "^${KEY}=" "$ENV_FILE"; then
  sed -i "s|^${KEY}=.*|${KEY}=${VALUE}|" "$ENV_FILE"
else
  echo "${KEY}=${VALUE}" >> "$ENV_FILE"
fi