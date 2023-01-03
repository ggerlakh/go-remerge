docker pull arangodb:latest
docker run --rm -e ARANGO_ROOT_PASSWORD=password -p 8529:8529 -d arangodb
