# ecommerce service

## Overview

ecommerce is a simple Go-based service to demostrate knowledge in golang.


## Setup

1. Install Go 1.23+
2. Run: `go run ./cmd/server`


## Running tests
To run the test locally use below command in cmd

    make test


## Environment bariables
    ENV=dev
    SERVER_PORT=5001
    AdminEmail=admin@example.com

    DB_CONNECTION_URL=postgres://ecommerce:ecommerce@2025@192.168.59.104:30007/ecommerce?sslmode=disable

    REDIS_ADDR=192.168.59.104:30008
    REDIS_PASSWORD=4e763fe4-4e54-45b0-ad10-1c449a0c24d0

    AFRICASTALKING_API_KEY=your_api_key
    AFRICASTALKING_USERNAME=sandbox
    AFRICASTALKING_SMTP_HOST=smtp.mailtrap.io
    AFRICASTALKING_SMTP_PORT=587
    AFRICASTALKING_SMTP_USER=user
    AFRICASTALKING_SMTP_PASS=pass

    EMAIL_FROM=admin@example.com
    EMAIL_PASSWORD=ds44534gvs
    EMAIL_SMTP_HOST=smtp.mailtrap.io
    EMAIL_SMTP_PORT=587
    ADMIN_EMAIL=admin@example.com

    OIDC_ISSUER="https://accounts.google.com"
    OIDC_AUDIENCE="your_audience"
    OIDC_JWKS_URL="https://www.googleapis.com/oauth2/v3/certs"


## Deployment
To deploy to cluster, include below below SECRETS to github

    TEST_CLUSTER - contains the kubernetes cluster yaml
    TEST_SERVICE_CONFIG - env variables to be included to variables
    GHCR_USERNAME - github username
    GHCR_TOKEN -  github token

Execute thw worflow called deploy

