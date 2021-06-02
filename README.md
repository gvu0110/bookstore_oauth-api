# bookstore_oauth-api

1. Start up Cassandra DB
Cassandra Docker run
```shell
docker run --name cassandra -d -p 9042:9042 cassandra:3.11.2
docker exec -it cassandra cqlsh -e "CREATE KEYSPACE oauth WITH REPLICATION = {'class':'SimpleStrategy', 'replication_factor':1}"
docker exec -it cassandra cqlsh -e "USE oauth; CREATE TABLE access_tokens(access_token varchar PRIMARY KEY, user_id bigint, client_id bigint, expires bigint);"
```
or
```shell
docker-compose up -d cassandra cassandra-load-keyspace
```
2. Run locally oauth-api app
```shell
export CASSANDRA_ADDRESS=localhost:9042
export USER_API_ENDPOINT=http://localhost:8081/users/login
go run main.go
```
3. Unit test coverage
```shell
go test -v -coverprofile cover.out .
go tool cover -html=cover.out -o cover.html
```