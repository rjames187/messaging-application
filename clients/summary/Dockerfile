FROM alpine as config-generator
RUN apk add --no-cache gettext
WORKDIR /app
COPY config.js.template .
CMD envsubst < config.js.template > /output/config.js

FROM nginx
COPY . /usr/share/nginx/html
ADD default.conf /etc/nginx/conf.d/default.conf
COPY --from=config-generator /output/config.js /usr/share/nginx/html/