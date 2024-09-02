FROM debian:stable-slim

COPY ./.bin/server ./bin/server
CMD [ "/bin/server" ]
