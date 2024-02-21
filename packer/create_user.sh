#!/bin/bash
echo "Updating the system..."
sudo dnf update -y
echo "Installing wget..."
sudo dnf install -y wget tar

echo "Creating the csye6225 group..."
sudo useradd -s /sbin/nologin -M csye6225