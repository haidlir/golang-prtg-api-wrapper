#!/bin/sh
sonar-scanner -X\
  -Dsonar.projectKey=golang-prtg-api-wrapper \
  -Dsonar.organization=haidlir-github \
  -Dsonar.sources=./prtg-api \
  -Dsonar.host.url=https://sonarcloud.io \
  -Dsonar.login=$LOGIN_SONARQUBE