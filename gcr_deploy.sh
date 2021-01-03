ENV_VARS=$(sed -E "s/^([^=]*)=(.*)$/--set-env-var '\1'=\2/" .env | tr '\n' ' ')
gcloud builds submit --tag gcr.io/cometplanning/comet-data-service
gcloud beta run deploy --image gcr.io/cometplanning/comet-data-service "${ENV_VARS}"
