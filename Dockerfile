FROM debian:stable-slim

COPY ./.bin/auth_serv ./bin/auth
CMD [ "/bin/auth" ]
