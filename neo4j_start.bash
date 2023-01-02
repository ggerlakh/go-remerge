#!bin/bash
docker pull neo4j:latest
docker run \
	--name testneo4j \
	--rm \
	-p7474:7474 -p7687:7687 \
	-d \
	-v $PWD/neo4j/data:/data \
	-v $PWD/neo4j/logs:/logs \
	-v $PWD/neo4j/import:/var/lib/neo4j/import \
	-v $PWD/neo4j/plugins:/plugins \
	--env NEO4J_AUTH=neo4j/neo4jdevops \
	neo4j:latest
