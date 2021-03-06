# Backend Deployment

## Initialize Firewall

ssh root@[server ip]

ufw allow OpenSSH

ufw enable

exit

## Create User
ssh root@[server ip]

adduser [user]

usermod -aG sudo [user]

rsync --archive --chown=[user]:[user] ~/.ssh /home/[user]

exit

## Install Golang Dependencies

cd ~

curl -O https://dl.google.com/go/go1.13.9.linux-amd64.tar.gz

sudo tar -xvf go1.13.9.linux-amd64.tar.gz -C /usr/local

mkdir -p $HOME/go/{bin,src}

nano ~/.profile

export GOPATH=$HOME/go

export PATH=$PATH:$GOPATH/bin

export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

. ~/.profile

## Pull Backend Project

cd $GOPATH/src

env GIT_TERMINAL_PROMPT=1 go get github.com/GalaSejahtera/GalaBackend

## Build Backend Project

cd github.com/GalaSejahtera/GalaBackend

go build cmd/server/main.go

## Initialize Service

sudo nano /lib/systemd/system/galasejahterabe.service

[Unit]

Description=Golang Service

ConditionPathExists=/home/lewis/go/src/github.com/GalaSejahtera/GalaBackend/main

After=network.target

[Service]

Restart=on-failure

RestartSec=10

startLimitIntervalSec=60

WorkingDirectory=/home/lewis/go/src/github.com/GalaSejahtera/GalaBackend

ExecStart=/home/lewis/go/src/github.com/GalaSejahtera/GalaBackend/main

[Install]

WantedBy=multi-user.target

## Input Environment Variables

systemctl edit galasejahterabe

[Service]

Environment="AUTH_ENABLED=true"

Environment="BACKEND_USER=ee00c51c-5ce2-461d-9526-862e2830a113"

Environment="MONGODB_URL=mongodb://localhost:27017"

Environment="ACCESS_SECRET=$2a$04$gKZuI64CEh7Jovqf/UTefO58PwT1YuGooEHG/W7oX8AYKFqWd9tQ6"

Environment="REFRESH_SECRET=$2a$04$Sy7aXxwd.auGV4diE8u3iueKpAS.vaY4XSA.Zji6pGOrrY8Ox3eHa"

Environment="EMAIL_ADDRESS=galasejahtera2020@gmail.com"

Environment="EMAIL_PASSWORD=GalaSejahtera2020!"

Environment="SMTP_SERVER_HOST=smtp.gmail.com"

Environment="SMTP_SERVER_PORT=587"

Environment="TEMPLATES_PATH=templates"

Environment="ADMIN_URL=http://localhost"

## Start Service

sudo service galasejahterabe start

## Restart Service

sudo service galasejahterabe restart

## Expose Port in Firewall

sudo ufw allow 10001/tcp

sudo ufw allow 10002/tcp

## View Service

sudo lsof -i -P -n | grep LISTEN

## .Profile Shortcut

navbe() {

        cd $HOME/go/src/github.com/GalaSejahtera/GalaBackend

}

deploybe() {

        cd $HOME/go/src/github.com/GalaSejahtera/GalaBackend

        go build cmd/server/main.go

        cd ~

        sudo service galasejahterabe restart

}

## MongoDB

sudo systemctl status mongodb

sudo apt install -y mongodb

sudo systemctl stop mongodb

sudo systemctl start mongodb

sudo systemctl restart mongodb

db.zones.createIndex( { "location" : "2dsphere" } )

db.users.createIndex( { "location" : "2dsphere" } )

db.authtokens.createIndex( { "ttl": 1 }, { expireAfterSeconds: 1 } )

db.activities.createIndex( { "ttl": 1 }, { expireAfterSeconds: 1 } )