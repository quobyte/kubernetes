#!/bin/bash

uname -a

echo registry=$QUOBYTE_REGISTRY > /etc/quobyte/host.cfg

if [[ $QUOBYTE_RPC_PORT ]]
then
  echo rpc.port=$QUOBYTE_RPC_PORT > /etc/quobyte/$QUOBYTE_SERVICE.cfg
fi

if [[ $QUOBYTE_HTTP_PORT ]]
then
  echo http.port=$QUOBYTE_HTTP_PORT >> /etc/quobyte/$QUOBYTE_SERVICE.cfg
fi

if [[ $QUOBYTE_WEBCONSOLE_PORT ]]
then
  echo webconsole.port=$QUOBYTE_WEBCONSOLE_PORT >> /etc/quobyte/$QUOBYTE_SERVICE.cfg
fi

echo test.device_dir=/devices >> /etc/quobyte/$QUOBYTE_SERVICE.cfg
echo public_ip=$HOST_IP >> /etc/quobyte/$QUOBYTE_SERVICE.cfg

SERVICE_UUID=$(uuidgen)
echo uuid=$SERVICE_UUID >> /etc/quobyte/$QUOBYTE_SERVICE.cfg

/etc/init.d/quobyte-$QUOBYTE_SERVICE start
/etc/init.d/quobyte-$QUOBYTE_SERVICE status

tee < /var/log/quobyte/$QUOBYTE_SERVICE.log

PID=$(cat /var/run/$QUOBYTE_SERVICE.run)

echo "Running Quobyte service $QUOBYTE_SERVICE $SERVICE_UUID in container as pid $PID"

while ps -p $PID > /dev/null; do sleep 1; done;
