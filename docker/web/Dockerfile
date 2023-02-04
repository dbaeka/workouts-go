FROM node:13.0-alpine3.11

RUN apk add yarn
ADD start.sh /
RUN chmod +x /start.sh

CMD ["/start.sh"]