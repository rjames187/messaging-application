### to run the server in a production VM copy and run the below commands

````docker pull rjames187/gateway:1.0

export TLSCERT=/etc/letsencrypt/live/api.rjames.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.rjames.me/privkey.pem

docker run -d \
--name gateway \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
rjames187/gateway:1.0```
````
