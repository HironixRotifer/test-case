#!/bin/bash

CONTAINER_NAME=api-gateway

if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    clear
    echo "Перезапуск контейнера: $CONTAINER_NAME"

    docker-compose stop $CONTAINER_NAME
    docker-compose rm $CONTAINER_NAME y
    docker-compose up -d --build $CONTAINER_NAME 
    
    echo "Контейнер $CONTAINER_NAME перезапущен."

    docker-compose logs -f $CONTAINER_NAME
else
    echo "Контейнер $CONTAINER_NAME не найден или не запущен."
fi