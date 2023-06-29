#!/bin/sh 


# exist(0) immediately 
set -e 

echo "run db migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_Source" -verbose up 

echo "Start the app"
exec "$@"

