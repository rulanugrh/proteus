#!/usr/bin/env bash

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

requirements() {
    check_os
}

build() {
    docker compose up -d db-product && docker compose up -d db-user && docker compose up -d db-order && docker compose up -d
}

help() {
    echo "Usage: $(basename "$0") [OPTIONS]
    Commands:
        requirements    To install requirements
        build           To build docker compose
        help            To show help command
    "
}

if [[ $1 =~ ^(requirements|build|help)$ ]]; then
    "$@"
else
    echo "Invalid command '$1'" >&2
    exit 1
fi