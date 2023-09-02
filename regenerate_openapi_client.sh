#!/bin/zsh
mkdir openapi
openapi-generator generate -i selfservice-api.yaml -g go --output ./openapi
cd ./openapi
go get github.com/stretchr/testify/assert
go get golang.org/x/net/context
find . -type f -exec sed -i '' 's|github.com/GIT_USER_ID/GIT_REPO_ID|github.com/centine/selfservice-api-client-go|g' {} +
