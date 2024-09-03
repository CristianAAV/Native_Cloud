#!/bin/bash

mkdir coverage
go test -v ./... -cover -covermode=atomic -coverprofile=coverage/coverage.out
go tool cover -func=coverage/coverage.out
go tool cover -html=coverage/coverage.out 

coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$coverage < 70" | bc -l) )); then
  echo "Test coverage is below 70%! The minimum required is 70%."
  exit 1
fi
echo "All tests passed and coverage is above 70%."