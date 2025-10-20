#!/bin/bash

cd "$(dirname "$0")"
echo ">> Start"

# Variables
NAME_CONTAINER="cnt_apigo_chatbot"
NAME_IMAGE="img_apigo_chatbot"
PORT_EXPOSE=3434
CONFIG_FOLDER="/config/chatbot/"
TIMEZONE="America/Guayaquil"

docker build -t "$NAME_IMAGE" -f Dockerfile .
echo ">> Compiled successfully"

# Si el contenedor existe, bórralo
if docker ps -a --format '{{.Names}}' | grep -q "^$NAME_CONTAINER$"; then
    echo ">> Removing the existing container"
    docker rm -f "$NAME_CONTAINER"
fi

docker run -d \
    --restart=always \
    --name "$NAME_CONTAINER" \
    -v "$CONFIG_FOLDER"config.json:/app/config.json \
    -e TZ="$TIMEZONE" \
    -p "$PORT_EXPOSE":8080 \
    "$NAME_IMAGE"

echo ">> Successfully service"
