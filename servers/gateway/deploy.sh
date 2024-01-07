# change the domain name on line 6 as needed
export TLSCERT=/etc/letsencrypt/live/api.rjames.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.rjames.me/privkey.pem
docker rm -f gateway
docker pull rjames187/gateway:1.0
docker run -d --name gateway -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY rjames187/gateway:1.0
exit