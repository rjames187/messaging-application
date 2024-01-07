# change the ip address in the ssh line as needed
./build.sh
docker push rjames187/summary-client:1.0
ssh -tt root@167.99.144.18 << END
docker rm -f summary-client
docker pull rjames187/summary-client:1.0
docker run -d --name summary-client -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=\$TLSCERT -e TLSKEY=\$TLSKEY rjames187/summary-client:1.0
exit
END