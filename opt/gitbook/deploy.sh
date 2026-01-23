#!/bin/bash

#declare REGISTRY=usernx
declare REGISTRY_REPO=usernx/p5
declare SRV=trade-master-docs
declare VER=1.0.0

docker build -t ${REGISTRY_REPO}:${VER} --build-arg SRV=${SRV} .
docker push ${REGISTRY_REPO}:${VER}