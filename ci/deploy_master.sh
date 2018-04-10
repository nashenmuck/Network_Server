#!/bin/bash
docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
docker tag nashenmuck/network:$TRAVIS_BRANCH nashenmuck/network:latest
docker tag nashenmuck/network:latest nashenmuck/network:latest-$(echo $TRAVIS_COMMIT | cut -c1-7)
docker push nashenmuck/network:latest
docker push nashenmuck/network:latest-$(echo $TRAVIS_COMMIT | cut -c1-7)

