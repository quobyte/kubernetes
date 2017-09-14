#!/usr/bin/env bash

echo "10.10.0.11 qb1" >> /etc/hosts
echo "10.10.0.12 qb2" >> /etc/hosts
echo "10.10.0.13 qb3" >> /etc/hosts
echo "10.10.0.14 qb4" >> /etc/hosts

yum install -yy epel-release
yum update
yum install ntp smartmontools

echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
sysctl -p /etc/sysctl.conf

