#!/usr/bin/env bash
tag=healthcheck-exporter:0.1.22

echo Building $tag
docker build --no-cache -t ziiot-docker.dp.nlmk.com/digital-plant/$tag . -f ./Dockerfile
docker push ziiot-docker.dp.nlmk.com/digital-plant/$tag