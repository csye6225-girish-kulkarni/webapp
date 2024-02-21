#!/bin/bash

sudo dnf install postgresql-server postgresql-contrib -y

sudo postgresql-setup initdb

sudo systemctl start postgresql

sudo systemctl enable postgresql

# Modify the PostgreSQL configuration file to listen on all addresses
sudo sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/g" /var/lib/pgsql/data/postgresql.conf

# Replace all occurrences of 'ident' with 'md5' in pg_hba.conf
sudo sed -i "s/ident/md5/g" /var/lib/pgsql/data/pg_hba.conf

sudo systemctl restart postgresql

echo "Creating PostgreSQL user and granting privileges..."
cd /tmp
sudo -u postgres psql -c "CREATE USER $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE postgres TO $POSTGRES_USER;"

# Install the uuid extension
sudo -u postgres psql -d postgres -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

echo "PostgreSQL installation and configuration complete."