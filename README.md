Code for Infrastructure as Code is Present in the following Repo
https://github.com/csye6225-girish-kulkarni/tf-gcp-infra


2) The Postgres Connection String is :

```
export POSTGRES_CONN_STR=postgres://girish:test@123@localhost:5432/postgres


```

```
export GOOS=linux
export GOARCH=amd64
```


## Command to build the go app
```
go build -o webapp main.go
```

## Transfer the binary to the VM 
```
scp webapp root@ipAddress
```

## script to install go on VM

```
#!/bin/bash
dnf install wget -y
# Download the Go binary archive
echo "Downloading Go 1.21.6..."
wget -q https://golang.org/dl/go1.21.6.linux-amd64.tar.gz

# Extract the archive to /usr/local
echo "Extracting Go 1.21.6..."
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# Set up Go environment variables
echo "Setting up Go environment variables..."
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# Verify the installation
echo "Verifying Go installation..."
go version



```
