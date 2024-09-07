#!/bin/bash

# Check if a migration name was provided
if [ -z "$1" ]; then
  echo "Usage: $0 \"migration name\""
  exit 1
fi

# Get the latest migration number
latest_migration=$(ls db/migrations | grep -E '^[0-9]{4}_.+\.up\.sql$' | sort | tail -n 1)
if [ -z "$latest_migration" ]; then
  number="0001"
else
  number=$(printf "%04d" $(( $(echo $latest_migration | cut -d'_' -f1) + 1 )))
fi

# Sluggify the migration name
name=$(echo "$1" | tr ' ' '_' | tr '[:upper:]' '[:lower:]')

# Create the migration files
up_file="db/migrations/${number}_${name}.up.sql"
down_file="db/migrations/${number}_${name}.down.sql"
touch "$up_file" "$down_file"

echo "Created migration files: $up_file and $down_file"
