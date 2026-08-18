#!/bin/bash

# Exit if any command fails
set -e

# Check if swagger.yaml exists
if [ ! -f "./docs/swagger.yaml" ]; then
  echo "ERROR! Could not find config file at path './docs/swagger.yaml'!"
  exit 1
fi

# Determine the current git branch
BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Set API, config, and gateway names based on the branch
if [ "$BRANCH" == "master" ]; then
  API_NAME="nebula-api"
  CONFIG_NAME="prod-config"
  GATEWAY_NAME="api-gateway"
else
  API_NAME="dev-nebula-api"
  CONFIG_NAME="dev-config"
  GATEWAY_NAME="dev-gateway"
fi

# Create temp config
echo "Creating temp config."
gcloud api-gateway api-configs create "${CONFIG_NAME}-temp" \
  --api="${API_NAME}" \
  --openapi-spec=./docs/swagger.yaml \
  --display-name="${CONFIG_NAME}-temp" \
  --quiet

# Migrate to temp config
echo "Migrating to temp config."
gcloud api-gateway gateways update "${GATEWAY_NAME}" \
  --location=us-central1 \
  --api="${API_NAME}" \
  --api-config="${CONFIG_NAME}-temp" \
  --quiet

# Delete old config
echo "Deleting old config."
gcloud api-gateway api-configs delete "${CONFIG_NAME}" \
  --api="${API_NAME}" \
  --quiet

# Create new config
echo "Creating new config."
gcloud api-gateway api-configs create "${CONFIG_NAME}" \
  --api="${API_NAME}" \
  --openapi-spec=./docs/swagger.yaml \
  --display-name="${CONFIG_NAME}" \
  --quiet

# Migrate to new config
echo "Migrating to new config."
gcloud api-gateway gateways update "${GATEWAY_NAME}" \
  --location=us-central1 \
  --api="${API_NAME}" \
  --api-config="${CONFIG_NAME}" \
  --quiet

# Delete temp config
echo "Deleting temp config."
gcloud api-gateway api-configs delete "${CONFIG_NAME}-temp" \
  --api="${API_NAME}" \
  --quiet

echo "Done!"
