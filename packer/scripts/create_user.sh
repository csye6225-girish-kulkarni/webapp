#!/bin/bash
echo "Updating the system..."
sudo dnf update -y
echo "Installing wget..."
sudo dnf install -y wget tar

echo "Creating the csye6225 group..."
sudo groupadd csye6225
sudo useradd -s /sbin/nologin -M -g csye6225 csye6225
# Creating the log directory and setting the permissions
sudo mkdir /var/log/webapp/
sudo chown -R csye6225:csye6225 /var/log/webapp/