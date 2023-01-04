#!bin/bash
docker pull neo4j:latest
docker run \
	--name neo4j \
	--rm \
	-p7474:7474 -p7687:7687 \
	-d \
	-v $PWD/neo4jdb/data:/data \
	-v $PWD/neo4jdb/logs:/logs \
	-v $PWD/neo4jdb/import:/var/lib/neo4j/import \
	-v $PWD/neo4jdb/plugins:/plugins \
	--env NEO4J_AUTH=neo4j/neo4jdevops \
	neo4j:latest
