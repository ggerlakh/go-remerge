#!bin/bash
docker pull arangodb:latest
docker run --name arangodb --rm -e ARANGO_ROOT_PASSWORD=password -p 8529:8529 -d arangodb
