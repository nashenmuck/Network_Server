#!/bin/bash

if [[ $TRAVIS_BRANCH == testing ]]
then
	export TAG=testing
fi
