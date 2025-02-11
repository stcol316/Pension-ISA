#!/bin/bash

#Note: Script to generate secret files
# Check dir exists
mkdir -p ../secrets

# Generate password and write to files
openssl rand -base64 32 > ../secrets/db_password.txt
echo "dev_db" > ../secrets/db_name.txt
echo "dev_user" > ../secrets/db_user.txt

# Permissions
chmod 600 ../secrets/*
