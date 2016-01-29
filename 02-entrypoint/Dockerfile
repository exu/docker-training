FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y golang

RUN echo "dupa" > /etc/dupa.conf


ENTRYPOINT ["go"]
CMD ["help"]
