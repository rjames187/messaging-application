./build.sh
docker push rjames187/summary-client:1.0
ssh -tt ec2-user@ec2-54-145-220-95.compute-1.amazonaws.com << END
sudo service docker start
docker rm -f summary-client
docker pull rjames187/summary-client:1.0
docker run -d --name summary-client -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=\$TLSCERT -e TLSKEY=\$TLSKEY rjames187/summary-client:1.0
exit
END