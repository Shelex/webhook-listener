# webhook-listener

service to listen for webhooks persisting to redis

## Start

```bash
# start redis container
    docker run --name my-redis -p 6379:6379 -d redis
```

```bash
# start service
    make start
```

Swagger will be available at localhost:8080/swagger/  
Websockets url: ws://localhost:8080/listen/{channel}

## How it works

- subscribes to http post requests to http://localhost:8080/api/{channel}
- publish message to go channel used as pubsub
- notification module read messages from channel and send them via websockets to subscribers
- repository module read messages from channel and store them in redis
- api endpoints available for get and delete operations for messages stored in redis for specific channel
- websocket subscription option available to check incoming messages
- every 3rd hour at 0 minute cron task checks storage and clears messages older than 3 days
