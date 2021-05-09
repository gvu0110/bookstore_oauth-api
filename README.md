# bookstore_oauth-api

Cassandra Docker run
```
docker run --name cassandra -d -p 9042:9042 cassandra:3.11.2
docker exec -it cassandra  cqlsh -e "CREATE KEYSPACE oauth WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1}"
docker exec -it cassandra cqlsh -e "USE oauth; CREATE TABLE access_tokens(access_token varchar PRIMARY KEY, user_id bigint, client_id bigint, expires bigint);"
```
or
```
docker-compose up
```