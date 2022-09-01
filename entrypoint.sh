#!/usr/bin/env bash

if [ ! -z $INPUT_ENVIRONMENT_USERNAME ];
then echo $INPUT_ENVIRONMENT_PASSWORD | docker login $INPUT_ENVIRONMENT_REGISTRY  \
-u $INPUT_ENVIRONMENT_USERNAME --password-stdin
fi


docker run -v "/var/run/docker.sock":"/var/run/docker.sock" -v $INPUT_WORKSPACE:/var/www --entrypoint=$INPUT_SHELL $INPUT_ENVIRONMENT_IMAGE -c "${INPUT_COMMAND//$'\n'/;}"
echo "::set-output name=result_code::$?"
echo "::set-output name=pipeline_artifacts::test"
# Generate build_id