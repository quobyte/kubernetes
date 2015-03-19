#!/bin/bash
# Run a Quobyte Storage service in a Docker container.
# The service will look for devices in /mnt.
# Use mount--bind to mount any Quobyte devices there.

# one of registry, metadata, data, api, webconsole
#QUOBYTE_SERVICE=data
#QUOBYTE_RPC_PORT=12345
#QUOBYTE_HTTP_PORT=12346
#QUOBYTE_WEBCONSOLE_PORT=12346
#QUOBYTE_REGISTRY=prod.corp.quobyte.com

# We run with bridge networking because host networking has a bug
# with su on current kernels, see
#  https://github.com/docker/docker/issues/5899


PORT_MAPPINGS=""

if [[ $QUOBYTE_WEBCONSOLE_PORT ]]
then
  PORT_MAPPINGS="-p $QUOBYTE_WEBCONSOLE_PORT:$QUOBYTE_WEBCONSOLE_PORT"
fi

if [[ $QUOBYTE_RPC_PORT ]]
then
  PORT_MAPPINGS="$PORT_MAPPINGS -p $QUOBYTE_RPC_PORT:$QUOBYTE_RPC_PORT -p $QUOBYTE_RPC_PORT:$QUOBYTE_RPC_PORT/udp"
fi

if [[ $QUOBYTE_HTTP_PORT ]]
then
  PORT_MAPPINGS="$PORT_MAPPINGS -p $QUOBYTE_HTTP_PORT:$QUOBYTE_HTTP_PORT "
fi

echo $PORT_MAPPINGS

sudo docker run --rm -i -t --privileged=true \
  -e QUOBYTE_SERVICE=$QUOBYTE_SERVICE \
  -e QUOBYTE_REGISTRY=$QUOBYTE_REGISTRY \
  -e QUOBYTE_RPC_PORT=$QUOBYTE_RPC_PORT \
  -e QUOBYTE_HTTP_PORT=$QUOBYTE_HTTP_PORT \
  -e QUOBYTE_WEBCONSOLE_PORT=$QUOBYTE_WEBCONSOLE_PORT \
  -e HOST_IP=$(dig +short $HOSTNAME) \
  $PORT_MAPPINGS \
  -h $(hostname -f) \
  -v /mnt/:/devices \
  quobyte-service

