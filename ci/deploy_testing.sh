#!/bin/bash
docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
docker push nashenmuck/network:testing
docker push nashenmuck/network:testing-$(echo $TRAVIS_COMMIT | cut -c1-7)

