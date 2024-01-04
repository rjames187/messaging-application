# use this script to run the container locally with self-signed certs

openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj "//CN=localhost" -keyout privkey.pem -out fullchain.pem

docker run -d \
--name gateway \
-p 443:443 \
-v /$(pwd):/etc/certs:ro \
-e TLSCERT=/etc/certs/fullchain.pem \
-e TLSKEY=/etc/certs/privkey.pem \
rjames187/gateway:1.0