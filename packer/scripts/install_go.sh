#!/bin/bash

# Install wget and tar to fetch and extract Go
sudo dnf install -y wget tar

GO_VERSION="1.21.6"
wget https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz

sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

# Clean up the downloaded tarball
rm go${GO_VERSION}.linux-amd64.tar.gz

# Set up Go environment variables for all users
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
echo 'export GOPATH=$HOME/go' | sudo tee -a /etc/profile
echo 'export GOBIN=$HOME/go/bin' | sudo tee -a /etc/profile

# Set up Go environment variables for the current user
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOBIN=$HOME/go/bin' >> ~/.bashrc

# Source the profile and bashrc to load the new environment variables
source /etc/profile
source ~/.bashrc

# Temporarily set PATH to include Go for the current script execution
export PATH=$PATH:/usr/local/go/bin

# Check the Go version
go version