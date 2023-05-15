# webhook-listener

Service to listen for webhooks, save them to redis, provide UI and API.

## Example

[UI](https://webhooks.shelex.dev/) and [Swagger](https://webhooks.shelex.dev/swagger/)

## Start

```bash
# start redis container
    docker run --name my-redis -p 6379:6379 -d redis
```

```bash
# start service
    make start
```

```bash
# start webpage
    make web-dev
```

Swagger will be available at localhost:8080/swagger/  
Websockets url: ws://localhost:8080/listen/{channel}  
Webpage url: localhost:8080/

## How it works

- accepts http post requests to `https://webhooks.shelex.dev/api/{channel}`
- publish message to redis pubsub "webhooks" channel
- notification module read messages by `webhooks` subscription and send them via websockets to subscribers
- repository module read messages by `webhooks` subscription and store them in redis
- api endpoints available for get and delete operations for messages stored in redis for specific channel
- websocket subscription option available to check incoming messages
- every 3rd hour at 0 minute cron task checks storage and clears messages older than 3 days


## How to set up VPS

Ubuntu 22 is used.

### Connect to ssh

IP address `111.111.111.11` as an example.

```bash
    ssh -i ssh/vultr root@111.111.111.11
```

### Install deps

 - 
    ```bash
        sudo apt update
    ```
 - redis: https://www.digitalocean.com/community/tutorials/how-to-install-and-secure-redis-on-ubuntu-22-04
 - nginx: https://www.digitalocean.com/community/tutorials/how-to-install-nginx-on-ubuntu-22-04


### Prepare web and go builds

Pre-requisite: prepare `.env` and `web/.env` config files with proper values.

```
    make web-build
    make build
```

### Upload to vps

```
    scp -r ./web/build/* root@111.111.111.11:/var/www/webhooks.shelex.dev/html/
    scp ./bin/* root@111.111.111.11:/usr/bin
    scp .env root@111.111.111.11:/usr/bin
```

### Handle backend service with systemd

    - 
    ```bash
        nano /etc/systemd/system/webhook-listener.service
    ```

    - paste content from `configs/webhook-listener.service`
        
    - reload daemon, start service, check status
        ```bash
            systemctl daemon-reload
            service webhook-listener start
            service webhook-listener status
        ```

### Basic setup for nginx and certbot

 - https://www.nginx.com/blog/using-free-ssltls-certificates-from-lets-encrypt-with-nginx/
 - `sudo certbot --nginx -d webhooks.shelex.dev -d www.webhooks.shelex.dev`

### Configure nginx

 - 
    ```bash
        sudo nano /etc/nginx/sites-enabled/webhooks.shelex.dev
    ```
 - paste content from `configs/webhooks.shelex.dev`
 - remove default site from `/etc/nginx/sites-available`
 - make sure ports are not conflicting with apache (`sudo update-rc.d apache2 disable`)