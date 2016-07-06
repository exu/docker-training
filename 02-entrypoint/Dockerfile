FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y golang

RUN echo "some_value" > /etc/dupa.conf


ENTRYPOINT ["go"]
CMD ["help"]
