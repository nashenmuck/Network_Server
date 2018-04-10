#!/bin/bash
if [ "$TRAVIS_BRANCH" = "master" ]; then
	bash ci/deploy_master.sh;
else if [ "$TRAVIS_BRANCH" = "testing" ]; then
	bash ci/deploy_testing.sh;
else
	echo "Not on any known branch or is a pull request, not deploying";
fi
fi
