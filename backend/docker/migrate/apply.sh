#!/bin/sh
set -eu

dbs="${POSTGRES_DB:-job_hunting_dev} job_hunting_test"

for db in $dbs; do
  for migration in /sql/migrations/*.sql; do
    [ -e "$migration" ] || continue
    echo "Applying $(basename "$migration") to $db"
    psql -v ON_ERROR_STOP=1 -h db -U "${POSTGRES_USER:-postgres}" -d "$db" -f "$migration"
  done
done
