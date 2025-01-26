#!/bin/bash

docker_postgresql_hostname=projectsprint-project3
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

echo "Success create database and users"
