#!/bin/sh

redis-server --demonize yes
while ! redis-cli ping; do
    echo "Waiting for redis to start..."
    sleep 1
done 

echo "Redis is ready. Starting Go application"

/app/myapp