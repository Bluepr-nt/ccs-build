#!/usr/bin/env bash
set -v
ENV_LIST=env.list
if [ ! -z $INPUT_ENVIRONMENT_USERNAME ];
then echo $INPUT_ENVIRONMENT_PASSWORD | docker login $INPUT_ENVIRONMENT_REGISTRY  \
-u $INPUT_ENVIRONMENT_USERNAME --password-stdin
fi

env | grep -E '(GITHUB|RUNNER|CI|INPUT|ARTIFACT_NAME)' > env.list
echo "CCS_SHORTSHA=${GITHUB_SHA::7}" >> env.list
CCS_BUILD_ID=`echo "${GITHUB_SHA::7}-$(date '+%Y%m%dT%H%M%S')"` >> env.list
echo "CCS_BUILD_ID=$CCS_BUILD_ID" >> env.list

mkdir pipeline_artifacts

docker run -v "/var/run/docker.sock":"/var/run/docker.sock" -v $GITHUB_WORKSPACE:$INPUT_ENVIRONMENT_WORKDIR \
  -v $PWD/pipeline_artifacts/:$INPUT_ENVIRONMENT_WORKDIR/$INPUT_PIPELINE_ARTIFACTS_PATH -w $INPUT_ENVIRONMENT_WORKDIR --entrypoint=$INPUT_SHELL --name $CCS_BUILD_ID \
  --env-file env.list $INPUT_ENVIRONMENT_IMAGE \
  -c "mkdir $INPUT_ENVIRONMENT_WORKDIR/$INPUT_PIPELINE_ARTIFACTS_PATH; ${INPUT_COMMAND//$'\n'/;};RESULT_CODE=$?"

echo "::set-output name=result_code::$?"

# for every folder in $INPUT_WORKSPACE create a tar.gz
# for every file in $INPUT_WORKSPACE upload them to artifacts store

LIST=`ls pipeline_artifacts`
echo "::set-output name=pipeline_artifacts::${LIST//$'\n'/, }"

# Generate build_id