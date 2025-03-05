@echo off
:: GCP sadly doesn't allow us to update the yaml config of an existing gateway config or do direct replacement of an existing config
:: instead, we make a temp config with the new yaml, deploy it, delete the original config, make a new one with the new yaml, migrate to that, and then delete the temp configs

echo Make sure you updated the "x-google-backend" config value before running this!
echo[
pause

IF NOT EXIST .\docs\swagger.yaml (
	echo ERROR! Could not find config file at path ".\docs\swagger.yaml"!
)

:: create temp gateway config
echo Creating temp config.
gcloud api-gateway api-configs create prod-config-temp --api=nebula-api --openapi-spec=./docs/swagger.yaml --display-name=prod-config-temp

echo Migrating to temp config.
:: migrate to temp config
gcloud api-gateway gateways update api-gateway --location=us-central1 --api nebula-api --api-config prod-config-temp

echo Deleting old prod config.
:: delete original config that is no longer in use
gcloud api-gateway api-configs delete prod-config --api=nebula-api

echo Creating new prod config.
:: create new config with original name -- same as temp config
gcloud api-gateway api-configs create prod-config --api=nebula-api --openapi-spec=./docs/swagger.yaml --display-name=prod-config

echo Migrating to new prod config.
:: migrate to new config
gcloud api-gateway gateways update api-gateway --location=us-central1 --api nebula-api --api-config prod-config

echo Deleting temp config.
:: delete temp config
gcloud api-gateway api-configs delete prod-config-temp --api=nebula-api

echo Done!