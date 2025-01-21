#!/bin/bash

docker_postgresql_hostname=projectsprint-project2
# TODO: get the password from .env variables
service_password="postgres"

# Down docker service for this projects alongside its data
# so it start fresh
docker compose down -v

sleep 1

docker compose up -d

sleep 2

echo "Docker Development Service Started"

echo "Start initialize database users for postgres"

# 1. Create necessary schema, users, and access for database access
# TODO: automate discovery of service and user creations

# Ceate schemas
docker exec -i $docker_postgresql_hostname \
  psql -U postgres -d projectsone -c "CREATE SCHEMA activities;
  CREATE SCHEMA logs;
  CREATE SCHEMA users;

  CREATE ROLE activities_user WITH LOGIN PASSWORD '$service_password';
  CREATE ROLE logs_user WITH LOGIN PASSWORD '$service_password';
  CREATE ROLE users_user WITH LOGIN PASSWORD '$service_password';

  GRANT ALL PRIVILEGES ON SCHEMA activities TO activities_user;
  GRANT ALL PRIVILEGES ON SCHEMA logs TO logs_user;
  GRANT ALL PRIVILEGES ON SCHEMA users TO users_user;
  "

echo "Success create database and users"
