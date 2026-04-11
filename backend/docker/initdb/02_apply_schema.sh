#!/bin/bash
set -e

# 開発DBとテストDBの両方にスキーマを適用
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /sql/schema.sql
psql -U "$POSTGRES_USER" -d "job_hunting_test" -f /sql/schema.sql

echo "Schema applied to $POSTGRES_DB and job_hunting_test"
