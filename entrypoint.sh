#!/bin/sh
set -e

mkdir -p /metasource_db

# If database is empty download it

if [ -z "$(ls -A /metasource_db)" ]; then
    echo "Database not found, downloading..."
    mkdir -p /tmp/metasource_db
    ./metasource-cli -location /tmp/metasource_db database
    cp -r /tmp/metasource_db/* /metasource_db/
    rm -rf /tmp/metasource_db
fi

echo "Starting service in dispense mode..."
./metasource-cli -location /metasource_db dispense

