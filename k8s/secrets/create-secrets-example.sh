#!/bin/bash

kubectl delete secret broadcast-api-secrets --ignore-not-found

kubectl create secret generic broadcast-api-secrets \
  --from-literal=DB_HOST="" \
  --from-literal=DB_PORT="5432" \
  --from-literal=DB_USER="" \
  --from-literal=DB_PASSWORD="" \
  --from-literal=DB_NAME="" \
  --from-literal=DB_URL="" \
  --from-literal=JWT_SECRET="" \
  --from-literal=SMTP_HOST="" \
  --from-literal=SMTP_PORT="587" \
  --from-literal=SMTP_USERNAME="" \
  --from-literal=SMTP_PASSWORD="" \
  --from-literal=SMTP_FROM_ADDR=""

echo "Secrets created successfully!"