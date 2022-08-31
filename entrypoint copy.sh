#!/usr/bin/env bash

if [ ! -z $INPUT_USERNAME ];
then echo $INPUT_PASSWORD | docker login $INPUT_REGISTRY -u $INPUT_USERNAME --password-stdin
fi


exec docker run -v "/var/run/docker.sock":"/var/run/docker.sock" -v $GITHUB_WORKSPACE:/var/www $INPUT_IMAGE /bin/bash -c "${INPUT_COMMAND//$'\n'/;}"