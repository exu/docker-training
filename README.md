# Docker training basics for dev worklows

## TODO Intro

1.  Docker basics
    -   VM vs containers

        [![img](res/docker-training-vm.png)](res/docker-training-vm.png)
        [![img](res/docker-training-containers.png)](res/docker-training-containers.png)

2.  Structure: Images vs containers

    Image: like class in programming
    Container: Instance

    We're storing, pushing, changing classes, but not runned instances
    (we need to rebuild)

1.  Containers Repository

    <https://hub.docker.com/>

## Running existing containers

run ubuntu with sth..

```
docker run ubuntu cat /etc/passwd
docker run ubuntu apt-get
docker run -it ubuntu /bin/bash
```

## Creating First container

-   Shell app

Add lines to `Dockerfile`
```
FROM ubuntu:latest
CMD date
```
next you'll need to run

```
docker build -t mybuntu .
docker run -r mybuntu
```

-   Change to webapp

## PHP App with mysql

-   First app without compose

## Same with docker compose



## Networking

## Binding ports

## Localhost integration (--net=host)

## Linking containers

networking

## playing with Docker-compose scale
