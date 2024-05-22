#!/bin/sh
docker build -t cdc-agent -f ./cdc_agent/Dockerfile .
docker tag cdc-agent antrad1978/cdc-agent:1.0
docker push antrad1978/cdc-agent:1.0