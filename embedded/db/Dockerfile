FROM ubuntu:latest

RUN groupadd -r mongodb && useradd -r -g mongodb mongodb

RUN apt-get update && apt-get install -y gnupg2
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10

RUN apt-get install -y mongodb

RUN mkdir -p /data/db

VOLUME ["/data/db"]

WORKDIR /data/db

EXPOSE 27017
ENTRYPOINT ["/usr/bin/mongod", "--smallfiles"]



