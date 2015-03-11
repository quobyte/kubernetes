#!/bin/bash

QUOBYTE_REPO_ID=7a5e33df70582afac43c4c5833d8114e

sed "s/QUOBYTE_REPO_ID/$QUOBYTE_REPO_ID/g" Dockerfile_release.templ >Dockerfile

docker build -t quobyte-service .
