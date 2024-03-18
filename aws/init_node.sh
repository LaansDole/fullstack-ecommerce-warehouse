#!/bin/bash

sudo yum update -y
sudo yum install git -y

# Clone the feature/aws branch of my repository
git clone -b feature/aws https://github.com/LaansDole/fullstack-ecommerce-warehouse.git

# Add the NodeSource repository - v18 (LTS)
curl -sL https://rpm.nodesource.com/setup_18.x | sudo bash -

sudo yum install nodejs npm -y

node -v
npm -v

# Install node version manager
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# Load nvm
source ~/.bashrc

# Use nvm to install the latest LTS version of NodeJS
nvm install --lts

# Verify that nodejs is installed
node -e "console.log('Running Node.js ' + process.version)"

### Database Setup ###
# https://dev.mysql.com/downloads/repo/yum/
sudo wget https://dev.mysql.com/get/mysql80-community-release-el9-5.noarch.rpm -y
sudo yum localinstall mysql80-community-release-el9-5.noarch.rpm -y
sudo yum install mysql-community-server -y
systemctl start mysqld.service

