@echo off
:: GCP sadly doesn't allow us to update the yaml config of an existing gateway config or do direct replacement of an existing config
:: instead, we make a temp config with the new yaml, deploy it, delete the original config, make a new one with the new yaml, migrate to that, and then delete the temp configs

IF NOT EXIST .\docs\swagger.yaml (
	echo ERROR! Could not find config file at path ".\docs\swagger.yaml"!
	exit /B 1
)

:: update prod or dev config depending on branch
FOR /F "usebackq delims=" %%i IN (`git branch --show-current`) do set BRANCH=%%i
IF "%BRANCH%"=="master" (
	set API_NAME=nebula-api
	set CONFIG_NAME=prod-config
	set GATEWAY_NAME=api-gateway
) ELSE (
	set API_NAME=dev-nebula-api
	set CONFIG_NAME=dev-config
	set GATEWAY_NAME=dev-gateway
)

:: create temp gateway config
echo Creating temp config.
call gcloud api-gateway api-configs create %CONFIG_NAME%-temp --api=%API_NAME% --openapi-spec=./docs/swagger.yaml --display-name=%CONFIG_NAME%-temp --quiet

echo Migrating to temp config.
:: migrate to temp config
call gcloud api-gateway gateways update %GATEWAY_NAME% --location=us-central1 --api %API_NAME% --api-config %CONFIG_NAME%-temp --quiet

echo Deleting old prod config.
:: delete original config that is no longer in use
call gcloud api-gateway api-configs delete %CONFIG_NAME% --api=%API_NAME% --quiet

echo Creating new prod config.
:: create new config with original name -- same as temp config
call gcloud api-gateway api-configs create %CONFIG_NAME% --api=%API_NAME% --openapi-spec=./docs/swagger.yaml --display-name=%CONFIG_NAME% --quiet

echo Migrating to new prod config.
:: migrate to new config
call gcloud api-gateway gateways update %GATEWAY_NAME% --location=us-central1 --api %API_NAME% --api-config %CONFIG_NAME% --quiet

echo Deleting temp config.
:: delete temp config
call gcloud api-gateway api-configs delete %CONFIG_NAME%-temp --api=%API_NAME% --quiet

echo Done!