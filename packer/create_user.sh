#!/bin/bash

echo "Installing wget..."
sudo dnf install wget -y

echo "Creating the csye6225 group..."
sudo useradd -s /sbin/nologin -M csye6225