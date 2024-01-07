# change the ip address in line 4 as needed
# change domain names in lines 6 and 7 as needed
./build.sh
docker push rjames187/gateway:1.0
ssh -tt root@142.93.0.40 << END
export TLSCERT=/etc/letsencrypt/live/api.rjames.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.rjames.me/privkey.pem
docker rm -f gateway
docker pull rjames187/gateway:1.0
docker run -d --name gateway -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=\$TLSCERT -e TLSKEY=\$TLSKEY rjames187/gateway:1.0
exit
END