#!/bin/bash

if [ "${TRAVIS_BRANCH}" == "master" ] && [ "${TRAVIS_PULL_REQUEST}" == "false" ]; then
    docker login -u "${DOCKER_IO_USERNAME}" -p "${DOCKER_IO_PASSWORD}" docker.io
    cd ./go/examples/knative-eventing-source
    make push-image
fi
