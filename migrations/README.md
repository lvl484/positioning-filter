# MIGRATIONS

Migrate with using golang/migrate default table/rules for user
In table must be created default user with name 'fuser' without any rules with own password

## environment variables

DBPORT -host port for database  
POSTGRES_USER -default owner for database  
POSTGRES_PASSWORD -password for owner  
POSTGRES_DB -default database name (db postgres will be also created)  
MIGRATIONS_PATH -path to folder with .sql migrations files  

## commands

docker run -v "$MIGRATIONS_PATH:/migrations"  --network host migrate/migrate -path=/migrations/ -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:$DB_POST/$POSTGRES_DB?sslmode=disable" up
