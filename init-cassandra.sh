#!/bin/bash

echo "CassandraSQL loading..."
sleep 30

echo "Creating Keyspace..."
docker exec -it shortfy-cassandra cqlsh -e "CREATE KEYSPACE IF NOT EXISTS shortfy WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};"

echo "Keyspace created!"
