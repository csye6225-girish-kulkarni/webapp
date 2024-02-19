#!/bin/bash

#sudo dnf update -y
#echo "updating dnf..."
#sudo dnf update -y
echo "Installing wget..."
sudo dnf install wget -y
# Create the csye6225 group
echo "Creating the csye6225 group..."
sudo useradd -s /sbin/nologin -M csye6225