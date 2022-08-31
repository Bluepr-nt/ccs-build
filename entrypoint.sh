#!/usr/bin/env bash

if [ ! -z $INPUT_ENVIRONMENT_USERNAME ];
then echo $INPUT_ENVIRONMENT_PASSWORD | docker login $INPUT_ENVIRONMENT_REGISTRY -u $INPUT_ENVIRONMENT_USERNAME --password-stdin
fi


exec docker run -v "/var/run/docker.sock":"/var/run/docker.sock" -v $INPUT_WORKSPACE:/var/www $INPUT_ENVIRONMENT_IMAGE "${INPUT_COMMAND//$'\n'/;}"