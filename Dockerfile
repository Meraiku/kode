FROM debian:stable

COPY ./.bin/auth_serv ./bin/auth

COPY ./certs/* /etc/ssl/certs/

CMD [ "/bin/auth" ]
