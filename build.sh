#!/usr/bin/env bash

command=$1

check_os() {
    . /etc/os-release
    case $ID in
        ubuntu)
        sudo apt install wget curl docker docker.io -y
        ;;
        arch)
        sudo pacman -S docker wget curl
        ;;
        debian)
        sudo apt install wget curl docker docker.io -y
        ;;
        centos)
        sudo yum install -y yum-utils && sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo && sudo yum install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
        ;;
    esac
}

if [[ $command == "requirements" ]]
then
    check_os
    exit 0
else if [[ $command == "build" ]]
then
    docker compose up -d db-product && docker compose up -d db-user && docker compose up -d
    exit
else
    echo "sorry invalid command"
fi
    