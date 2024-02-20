#!/bin/bash
# Disabling SELinux
echo "SELINUX=permissive" | sudo tee /etc/selinux/config
sudo setenforce 0
sudo mv /tmp/webapp.service /etc/systemd/system/webapp.service

port_number=8080
# Add a firewall rule to allow traffic on the port
sudo firewall-cmd --permanent --add-port=${port_number}/tcp

# Reload the firewall to apply the changes
sudo firewall-cmd --reload

sudo systemctl daemon-reload
sudo systemctl enable webapp.service
sudo systemctl start webapp.service