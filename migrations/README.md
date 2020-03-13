# MIGRATIONS

Service creates container with PostgreSQL and migrate there default table/rules for user  

## environment variables

DBPORT -host port for database  
POSTGRES_USER -default owner for database  
POSTGRES_PASSWORD -password for owner  
POSTGRES_DB -default database name (db postgres will be also created)  
MIGRATIONS_PATH -path to folder with .sql migrations files  

## commands

docker-compose up -d   runs docker container with postgreSQL database  

sudo docker run -v "$MIGRATIONS_PATH:/migrations"  --network host migrate/migrate -path=/migrations/ -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$DB_POST/$POSTGRES_DB?sslmode=disable" up
