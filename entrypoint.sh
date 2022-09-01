#!/usr/bin/env bash

if [ ! -z $INPUT_ENVIRONMENT_USERNAME ];
then echo $INPUT_ENVIRONMENT_PASSWORD | docker login $INPUT_ENVIRONMENT_REGISTRY  \
-u $INPUT_ENVIRONMENT_USERNAME --password-stdin
fi


docker run -v "/var/run/docker.sock":"/var/run/docker.sock" -v $INPUT_WORKSPACE:$INPUT_ENVIRONMENT_WORKDIR -w $INPUT_ENVIRONMENT_WORKDIR --entrypoint=$INPUT_SHELL $INPUT_ENVIRONMENT_IMAGE -c "mkdir /var/www/${INPUT_COMMAND//$'\n'/;};RESULT_CODE=$?"
echo "::set-output name=result_code::$?"
LIST=`ls $INPUT_WORKSPACE`
echo "::set-output name=pipeline_artifacts::${LIST//$'\n'/, }"

# Generate build_id