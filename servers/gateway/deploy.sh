./build.sh
docker push rjames187/gateway:1.0
ssh -tt ec2-user@ec2-54-87-136-143.compute-1.amazonaws.com << END
sudo service docker start
docker rm -f gateway
docker pull rjames187/gateway:1.0
docker run -d --name gateway -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro -e TLSCERT=\$TLSCERT -e TLSKEY=\$TLSKEY rjames187/gateway:1.0
exit
END