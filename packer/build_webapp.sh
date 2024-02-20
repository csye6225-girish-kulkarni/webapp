#!/bin/bash


sudo mkdir -p /usr/bin/

sudo mv /tmp/webapp /usr/bin/
# Change the ownership of the webapp file to csye6225
sudo chown csye6225:csye6225 /usr/bin/webapp
# Change the permissions of the webapp file to make it executable
sudo chmod +x /usr/bin/webapp