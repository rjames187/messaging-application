# change the domain name on line 6 as needed
export TLSCERT=/etc/letsencrypt/live/rjames.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/rjames.me/privkey.pem
docker rm -f summary-client
docker pull rjames187/summary-client:1.0
docker run -d --name summary-client -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY rjames187/summary-client:1.0
exit