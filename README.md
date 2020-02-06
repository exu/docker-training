# What's Docker 

## Docker architecture

[![img](res/docker-architecture.svg)](res/docker-architecture.svg)

Docker parts: 
- *Docker daemon* - manages docker objects
- *Docker client* - cli for interacting with docker
- *Docker registries* - stores docker images
- *Docker objects* - when interacting with Docker you're managing objects. There are several types of objects
  - *Images* - template with instructions how to create docker container
  - *Containers* - is a runnable instance of an _Image_
  - *networks*, 
  - *volumes*, 
  - *plugins*, 

## Why Docker? It's all about standards and resource utilisation.

-   VM vs containers

[![img](res/docker-training-vm.png)](res/docker-training-vm.png)

[![img](res/docker-training-containers.png)](res/docker-training-containers.png)







# Hands on Docker examples

```
                    ##        .
              ## ## ##       ==
           ## ## ## ##      ===
       /""""""""""""""""\___/ ===
  ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~
       \______ o          __/
         \    \        __/ 
          \____\______/
```

## Docker Repositories

There is huge containers repository with official and unofficial images,
here you can find many of prebuild container images with different linux
distributions, web servers, applications, databases, production ready
systems and many many more.

<https://hub.docker.com/>

now let's run some containers

### Running existing containers

But how containers work? Let's start with really simple

example:

run ubuntu with sth..

```
docker run ubuntu cat /etc/passwd
docker run ubuntu apt-get
docker run -it ubuntu /bin/bash
```

What's happen when you're running docker containers?
- first docker will check if there is image which you want to run available
- if it is not it'll check if there are parts of images available 
- if not it downloads all needed parts of image 
- docker starts your image with given `cmd` (last passed parameter in examples
  above)


## Simple first run [Dockerfile](001-simple-dockerfile/Dockerfile)

Running images with command line is good when you want to check or debug
something quickly, but the main purpose of Docker is to make your application 
to be immutable with all dependencies (the system libs too)

To define such image you'll need some DSL. `Dockerfile` is such DSL in Docker.

Simplest Dockerfile could look like this: 

first create new Dockerfile `touch Dockerfile`

Add lines to `Dockerfile`
```
FROM ubuntu:latest
CMD date
```
next you'll need to build your image from Dockerfile and run built image.

```
docker build -t cmd .
docker run cmd
```

When we pass additional parameters We override CMD section.

```
docker run cmd ls -la
```

- Building with `docker build -t TAG_NAME`.
- Running  with `docker run TAG_NAME`


## Entrypoints [Dockerfile](002-entrypoint/Dockerfile)

Default entrypoing in docker is `/bin/sh -c` which simply runs command passed to `CMD` instruction


We can set entrypoint for our app (default is `/bin/sh -c`)
`CMD` will be appended.

```
FROM ubuntu:latest

ENTRYPOINT ["date", "-R"]
CMD ["-u"]
```

you can override command passing additional parameters after run:

```
docker build -t ep .
docker run ep --date='@1417400000'
```


## Inserting editor inside docker [Dockerfile](003-editor/Dockerfile)
 
You can run almost all applications / services in Docker. But sometimes you'll 
need to attach TTY to your docker container 

- build with `docker build -t editor` 
- run with `docker run -it editor`

## Docker images (ubuntu vs alpine vs debian vs scratch)

You need be careful in your choice of base docker image - they can be huge, 
and you for sure don't want to pass so big images through your network.

Please build image based on `ubuntu` next run one based on `alpine` images

``` Dockerfile
FROM alpine:3

RUN apk add curl

CMD curl ifconfig.co
```

``` Dockerfile
FROM ubuntu:latest

RUN apt update          # ubuntu doesn't update package urls by default
RUN apt install -y curl

CMD curl ifconfig.co
```


``` sh
$ docker build -t size-ubuntu -f Dockerfile.ubuntu .
$ docker build -t size-alpine -f Dockerfile.alpine .
```

Next check `docker images` command to check image size:

``` sh
$ docker images | head -n 5
```



## Long running processes in Docker

Recently our all examples was for running single run commands which ends their
life after execution. But often we'll be dealing with images which have some
long running processes. 

For this example we'll lock our Docker Image with some simple `sleep` command

``` Dockerfile
FROM alpine:3 

CMD ["tail", "-f", "/dev/stdout"] # tail -f will output if new lines will be
                                  # available and it'll lock our process
```

``` sh
$ docker build -t tail . 
$ docker run -it tail
```

Now we can go to new Terminal and check what's going on with that process.

``` sh
$ docker ps

# command output
CONTAINER ID        IMAGE               COMMAND                 CREATED             STATUS              PORTS               NAMES
72532d002a16        long                "tail -f /dev/stdout"   7 seconds ago       Up 6 seconds                            elegant_kilby
01029e402dde        ubuntu              "/bin/sh"               30 minutes ago      Up 30 minutes                           optimistic_morse

```

as we can se on my machine example output shows that there are two images:
- first was created 7 seconds ago
- second was run 30 minutes ago
- please notice that docker generetes random `NAME` for each container 

now we can check whats going on our run container

``` sh
#              attach terminal from our container
#               /
$ docker exec -it 72532d002a16 /bin/sh
#                    /            \ 
#               conainer ID      command 
#               or name
```

`docker exec` runs command inside working container
we're running here shell inside container (`-it` is needed to attach our shell
to docker shell)

## Named conainers

Running `docker exec` with conintaer ID passed (or random name) could be quite inconvinient - but there
is a nice option when running containers `--name` which sets a custom name for
container. 

``` sh
$ docker run --rm --name tailer tail
```

``` sh
$ docker ps 
CONTAINER ID        IMAGE               COMMAND                 CREATED             STATUS              PORTS               NAMES
8954e7811347        tail                "tail -f /dev/stdout"   6 seconds ago       Up 5 seconds                            tailer

```

as we can see that container `NAME` is set, now we can use it instead of ID's

e.g. 

``` sh
$ docker exec -it tailer /bin/sh
```

## Docker containers are immutable 

File system in docker containers is temporary by default (exception here are data volumes). 
When you stop container and start again all data will be lost same will happen with new instance of class in OOP. Persistance here is made on building process.

If you change something in your docker images it'll simply lost after docker container will be reloaded. 

``` sh
echo "some file" >> some_file.txt
echo "content" >> some_file.txt
echo "and more content" >> some_file.txt
```

``` Dockerfile
FROM alpine:3
COPY some_file.txt /
CMD tail -f /some_file.txt
```

``` sh
$ docker build -t immutable .
$ docker run -it --name im1 --rm immutable
```

next you can modify content of this container from other terminal session

``` sh
$ docker exec -it im1 /bin/sh
```

and add some files to our `some_file.txt` inside container

``` sh
$ echo "another line" >> /some_file.txt
$ echo "another line" >> /some_file.txt
$ echo "another line" >> /some_file.txt
$ echo "another line" >> /some_file.txt
$ echo "another line" >> /some_file.txt
```

as you can see on first terminal our file got new lines - file was modified

Now let's restart our container

``` sh
$ docker kill im1
$ docker run --rm --name im1 immutable
```

Let's get to it's shell again:
``` sh
$ docker exec -it im1 /bin/sh
```


``` sh
$ cat /some_file.txt
some file
content
and more content
```
as we can see after container restart file is not changed

Example - [Dockerfile](060-immutable-images/Dockerfile)

## Docker persistance

If you want to persist your data you'll need to use volumes or mounts we'll look at volumes first. 
It's like attaching new disk to your PC. You can attach multiple volumes to multiple directories.

![mounts](res/types-of-mounts-volume.png "mounts")

There are two types of mounts
- bind mount
- volume


Volumes are the *preferred mechanism for persisting data* generated by and used by
Docker containers. While bind mounts are dependent on the directory structure of 
the host machine, volumes are completely managed by Docker.

### Creating new volumes

You can explicitly create new named volumes: 

- craete new volume
``` sh
$ docker volume create myvol

myvol
```

- list created volumes
``` sh
$ docker volume ls 

DRIVER              VOLUME NAME
local               myvol                                                                                                                             
```

- inspect volume
``` sh
$ docker volume inspect myvol

[                                                                                                                                                     
    {                                                                                                                                                 
        "CreatedAt": "2020-02-01T08:03:58+01:00",                                                                                                     
        "Driver": "local",                                                                                                                            
        "Labels": {},                                                                                                                                 
        "Mountpoint": "/var/snap/docker/common/var-lib-docker/volumes/myvol/_data",                                                                   
        "Name": "myvol",
        "Options": {},
        "Scope": "local"
    }
]
```

- rm created volume
``` sh
$ docker volume rm myvol 

myvol
```


### Attaching volumes to conainer

``` sh
$ docker volume create vol1
```


``` sh
$ docker run -d \
  --name devtest \
  --mount source=vol1,target=/app \
  nginx:latest
```

Next lets run another container which will use our volume.

``` sh
$ docker run -it --mount source=vol1,destination=/mymount ubuntu /bin/sh
```

#### Some tips form docker site:
- Named volues are a lot easier to use and backup
- Originally, the -v or --volume flag was used for standalone containers and the --mount flag was used for swarm services. However, starting with Docker 17.06, you can also use --mount with standalone containers. In general, --mount is more explicit and verbose. The biggest difference is that the -v syntax combines all the options together in one field, while the --mount syntax separates them. Here is a comparison of the syntax for each flag.
- New users should try --mount syntax which is simpler than --volume syntax.


### Bind mounts 

The file or directory is referenced by its **full or relative path on the host
machine**. 

By contrast, when you use a volume, a new directory is created within
Docker’s storage directory on the host machine, and Docker manages that
directory’s contents.

To bind mount you can use `-v` paramter, values are separated by `:`
In the case of bind mounts, 
- the first field is the **path** to the file or directory on the **host machine**.
- The second field is the **path** where the file or directory is **mounted in the container**.
- The third optional is and is comma separated **list of options**.


to bind mount directory from local filesystem use following command: 

``` sh
#                   need to be full path
#                       /
$ docker run -it -v $(pwd)/dirToMount:/whereToMountInContainer ubuntu
#                              /           \
#                            local          container 
```

local dir will bin bound to container
volumes are often used to quickly run your apps with local configuration or to
debug changes


### Drivers 

There are many drivers for volumes, you can even use S3 or GCS
- https://docs.docker.com/registry/storage-drivers/gcs/


## `Dockerfile` reference 

We've learned several Dockerfile instructions in previous examples. But
Dockerfile has more options to use. 

### Build context

When you issue a docker build command, the **current working directory** is called
the **build context**. By default, the Dockerfile is assumed to be located here, but
you can specify a different location with the file flag (`-f`). 

Regardless of where the Dockerfile actually lives, all recursive contents of files and
directories in the current directory are sent to the Docker daemon as the build
context.

Do if your project is quite big it could be quite slow when you're running
`docker build` command. 

``` sh
$ docker build .
Sending build context to Docker daemon  6.51 MB
```

How to limit your build context? Simply use `.dockerignore` file, it's working like
.gitignore but for Docker context - you can easily limit not needed files and
directories like caches, intermediate build files or temporary directories.



## Dockerfile commands 

Docker file is based on instructions (you've seen them above of this document)

``` Dockerfile
# Comment
INSTRUCTION arguments
```

There are many instructions: 

### `FROM`

e.g. `FROM ubuntu:tag`, `FROM ubuntu:20.04 as newUbuntu`, `FROM ubuntu@j2j43r0924i3ir0234r092i43r`

You can use arguments with FROM: 
``` Dockerfile
ARG  CODE_VERSION=latest
FROM base:${CODE_VERSION}
```


### `RUN`
RUN has 2 forms:

``` Dockerfile
RUN <command> #(shell form, the command is run in a shell, which by default is /bin/sh -c on Linux or cmd /S /C on Windows)
RUN ["executable", "param1", "param2"] #(exec form)
```

**Layering** RUN instructions and generating commits conforms to the core concepts
of Docker where commits are **cheap** and containers can be created from any point
in an image’s history, much like source control.


### `CMD`
The CMD instruction has three forms:

``` Dockerfile
CMD ["executable","param1","param2"] #(exec form, this is the preferred form)
CMD ["param1","param2"] #(as default parameters to ENTRYPOINT)
CMD command param1 param2 #(shell form)
```

**The main purpose of a CMD is to provide defaults for an executing container**.
These defaults can include an executable, or they can omit the executable, in
which case you must specify an ENTRYPOINT instruction as well.


### `LABEL`

The LABEL instruction adds metadata to an image. A LABEL is a key-value pair. To
include spaces within a LABEL value, use quotes and backslashes as you would in
command-line parsing. A few usage examples:

``` Dockerfile
LABEL maintainer="jacek.wysocki@gmail.com"
LABEL "com.doctrime.version"="1.4.55-build-153"
LABEL com.example.label-with-value="foo"
LABEL version="1.0"
LABEL description="New basket service \
will replace part of monolith after implementing currency panel."
```

### `EXPOSE`

The EXPOSE instruction informs Docker that the container listens on the
specified network ports at runtime. You can specify whether the port listens on
TCP or UDP, and the default is TCP if the protocol is not specified

The EXPOSE instruction does not actually publish the port. It functions as a
type of documentation between the person who builds the image and the person who
runs the container, about which ports are intended to be published

``` Dockerfile
EXPOSE 8080
```

### `ENV`

The ENV instruction sets the environment variable <key> to the value <value>. 



``` Dockerfile
ENV myName John Doe
ENV myDog Rex The Dog
ENV myCat fluffy
```

The environment variables set using ENV will persist when a container is run
from the resulting image. You can view the values using docker inspect, and
change them using `docker run --env <key>=<value>`


You can replace defined envs inline details about rules can br found in docs:
https://docs.docker.com/engine/reference/builder/#environment-replacement


### `ADD`

The ADD instruction copies new files, directories or remote file URLs from <src>
and adds them to the filesystem of the image at the path <dest>.


ADD has two forms:

``` Dockerfile
ADD [--chown=<user>:<group>] <src>... <dest>
ADD [--chown=<user>:<group>] ["<src>",... "<dest>"] # (this form is required for paths containing whitespace)
```

Some examples: 

``` Dockerfile
ADD hom* /mydir/        # adds all files starting with "hom"
ADD hom?.txt /mydir/    # ? is replaced with any single character, e.g., "home.txt"

ADD test relativeDir/          # adds "test" to `WORKDIR`/relativeDir/
ADD test /absoluteDir/         # adds "test" to /absoluteDir/

ADD --chown=55:mygroup files* /somedir/
ADD --chown=bin files* /somedir/
ADD --chown=1 files* /somedir/
ADD --chown=10:11 files* /somedir/

ADD http://sources.file/somefile.txt /somedir/
ADD superarchive.tar.gz /somedir/
```

There are some limitations about add `src` and `dest` which you can find in
docs: https://docs.docker.com/engine/reference/builder/#add


### `COPY`

COPY is very similiar to add but allow to insert into container only files which
are in context. 

### `ENTRYPOINT`

``` Dockerfile
ENTRYPOINT ["executable", "param1", "param2"] #(exec form, preferred)
ENTRYPOINT command param1 param2 #(shell form)
```


Example:
``` Dockerfile
FROM ubuntu
ENTRYPOINT ["top", "-b"]
CMD ["-c"]
```

Both `CMD` and `ENTRYPOINT` instructions define what command gets executed when
running a container. There are few rules that describe their co-operation.

- `Dockerfile` should specify at least one of `CMD` or `ENTRYPOINT` commands.
- `ENTRYPOINT` should be defined when using the container as an executable.
- `CMD` should be used as a way of defining default arguments for an `ENTRYPOINT` command or for executing an ad-hoc command in a container.
- `CMD` will be overridden when running the container with alternative arguments.

### `VOLUME`

The `VOLUME` instruction creates a mount point with the specified name and marks
it as holding externally mounted volumes from native host or other containers.
The value can be a JSON array, `VOLUME ["/var/log/"]`, or a plain string with
multiple arguments, such as `VOLUME /var/log` or `VOLUME /var/log /var/db`.

For more information/examples and mounting instructions via the Docker client,
refer to [Share Directories via Volumes documentation](https://docs.docker.com/engine/tutorials/dockervolumes/).


``` Dockerfile
FROM ubuntu
RUN mkdir /myvol
RUN echo "hello world" > /myvol/greeting
VOLUME /myvol
```

This Dockerfile results in an image that causes docker run to create a new mount
point at `/myvol` and copy the greeting file into the newly created volume.


### `USER`

``` Dockerfile
USER <user>[:<group>] or
USER <UID>[:<GID>]
```

### `WORKDIR`
  

The `WORKDIR` instruction sets the working directory for any `RUN`, `CMD`, `ENTRYPOINT`,
`COPY` and `ADD` instructions that follow it in the Dockerfile. If the `WORKDIR`
doesn’t exist, it will be created even if it’s not used in any subsequent
Dockerfile instruction.

``` Dockerfile
ENV DIRPATH /path
WORKDIR $DIRPATH/$DIRNAME
RUN pwd
WORKDIR /someroot
RUN ls -la 
WORKDIR /
```

### `ARG`

The `ARG` instruction defines a variable that users can pass at build-time to the
builder with the docker build command using the `--build-arg <varname>=<value>`
flag. 

``` Dockerfile
FROM busybox
ARG user1=jacekwysocki
ARG gitHash=ii3e09i329e23
ARG buildno=19320
RUN echo "${buildno}"
```

### `ONBUILD`

The `ONBUILD` instruction adds to the image a trigger instruction to be executed
at a later time, when the **image is used as the base for another build**. 

The trigger will be executed in the context of the downstream build, as if it
had been inserted immediately after the FROM instruction in the downstream Dockerfile.

``` Dockerfile
ONBUILD RUN /usr/local/bin/python-build --dir /app/src
```


### `STOPSIGNAL`

The STOPSIGNAL instruction sets the system call signal that will be sent to the
container to exit. This signal can be a valid unsigned number that matches a
position in the kernel’s syscall table, for instance 9, or a signal name in the
format SIGNAME, for instance SIGKILL.

### `HEALTHCHECK`

``` Dockerfile
HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost/ || exit 1
```

To help debug failing probes, any output text (UTF-8 encoded) that the command
writes on stdout or stderr will be **stored** in the **health status** and can be
queried with docker inspect. Such output should be kept short (only the first
4096 bytes are stored currently).





# Creating applications

```
                    ##        .
              ## ## ##       ==
           ## ## ## ##      ===
       /""""""""""""""""\___/ ===
  ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~
       \______ o          __/
         \    \        __/ 
          \____\______/
```


## Creating Simple PHP Web application

## Getting our container IP address

Ok we've run our web app, It's working but how to show it in browser? We'll need
our container IP address:

```
docker inspect
```

We can filter inspected data with format parameter: 

```
docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $INSTANCE_ID
```

result:

```
172.17.0.4
```

http://172.17.0.4:8080/allaallalal




### Golang based (app server)

### PHP based application

## Exposing application





    
# Simple web app

## Building inside container

First we build our statically linked `app.go` file:

```
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi %s!!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

With `go build -o app` You'll receive `app` binary which will be used in our
docker container. We'll add file inside our container on build

```
FROM ubuntu:latest
ADD ./app /srv/app
CMD /srv/app
```

```
docker build -t swa .
```

We run our container with name

```
docker run --name=gogo swa
```


Now when we propagate container it'll be freezed inside it.

# Data volumes

I've said on beggining that docker has temporary file system by default.
But what when you want to store some informations after container will
be removed?


# PHP App with mysql

-   First app without compose

# Binding ports

- `-P` bind all ports to local machine high ports (from ephemeral port range which typically ranges from 32768 to 61000)
- `-p 5000` binds port 5000 from container to high port
- `-p 4900:5000` binds port 5000 from container to 4900 port on local machine


# Connection containers - Networking

In docker 1.8 and below links between containers was used, You'll need
to explicitly set link between two containers.

From 1.9 valid connection between containers is make with use of networking
First create network:

    docker network create training1

Then you'll need to pass `--net=training1` to `docker run` command


## Docker compose

When you run `docker-compose --x-networking up` in `myapp` directory, the following happens:

- A network called myapp is created.
- A container is created using web’s configuration. It joins the network myapp under the name myapp_web_1.
- A container is created using db’s configuration. It joins the network myapp under the name myapp_db_1.


# Localhost integration (--net=host)


# Pushing your image to public

Get your image id

```
docker images | less
```


Make tag

```
                              docker hub username
                            /      repo name
                           /     /    tag name
                          /     /   /
docker tag e701985f4a8c ex00/emacs:v1
               \
                image id


```

Login to hub and publish your image

```
docker login --username=ex00 --email=jacek.wysocki@gmail.com

docker push
```


# Your images are too big (could be > 700MB per image)

You can use micro base images:

- `phusion/baseimage` ~6MB with apt-get based on ubuntu lts. with init process
- `gliderlabs/alpine` ~ 5MB base image


# playing with Docker-compose scale
  
# Scaling single-core apps with static docker config

- http://blog.hypriot.com/post/docker-compose-nodejs-haproxy/



# Useful images

* Consul
  - https://hub.docker.com/r/voxxit/consul/
  - old: https://hub.docker.com/r/progrium/consul/
* Registrator
  - http://gliderlabs.com/registrator/latest/user/quickstart/
* PHP
* Webservers / Load balancers




# Scheduling


## Docker Swarm Mode! from version 1.12 available as part of docker

### What is Swarm

Docker Swarm is cluster manager and container scheduler, it is responsible
for puting given container instance on apropriate cluster node.

Imagine some web application...

Imagine that we have 10 nodes cluster, to simplfy let's say that each node
have two processor cores and 4GB of RAM.

Imagine also that your application have some kind of long running worker which
is responsible for sending emails

You have ~ 1 000 000 users

Almost all users are using your app on monday so your app servers are overloaded

Your users are sending emails each friday for their friends inside applications so your
workers are screaming for more processor power

Your job is to minimize cost, you can do it manually, each friday put several worker processes more
to each workers node where they can fight for resources with other.

On monday You'll check if any of workers instances are in idle mode and you can kill them (probably
with use of some kind of admin panel. And spawn more www processes which can fight for resources.

When you add more nodes it'll be more complicated.


Here comes swarm (or kubernetes) to help Us

Swarm organize your machine in cluster, is doing load balancing of your apps for You. You are able to
put your containerized app to your cluster and swarm will run it where there are resources
available.

Swarm

Terms:
- A swarm is a cluster of Docker Engines where you deploy services
- A node is an instance of the Docker Engine participating in the swarm.
- A service is the definition of the tasks to execute on the worker nodes. It is the central structure of the swarm system and the primary root of user interaction with the swarm.
- A task carries a Docker container and the commands to run inside the container


- Cluster management integrated with Docker Engine: Use the Docker Engine CLI to create a Swarm of Docker Engines where you can deploy application services. You don’t need additional orchestration software to create or manage a Swarm.

- Decentralized design: Instead of handling differentiation between node roles at deployment time, the Docker Engine handles any specialization at runtime. You can deploy both kinds of nodes, managers and workers, using the Docker Engine. This means you can build an entire Swarm from a single disk image.

- Declarative service model: Docker Engine uses a declarative approach to let you define the desired state of the various services in your application stack. For example, you might describe an application comprised of a web front end service with message queueing services and a database backend.

- Scaling: For each service, you can declare the number of tasks you want to run. When you scale up or down, the swarm manager automatically adapts by adding or removing tasks to maintain the desired state.

- Desired state reconciliation: The swarm manager node constantly monitors the cluster state and reconciles any differences between the actual state your expressed desired state. For example, if you set up a service to run 10 replicas of a container, and a worker machine hosting two of those replicas crashes, the manager will create two new replicas to replace the ones that crashed. The swarm manager assigns the new replicas to workers that are running and available.

- Multi-host networking: You can specify an overlay network for your services. The swarm manager automatically assigns addresses to the containers on the overlay network when it initializes or updates the application.

- Service discovery: Swarm manager nodes assign each service in the swarm a unique DNS name and load balances running containers. You can query every container running in the swarm through a DNS server embedded in the swarm.

- Load balancing: You can expose the ports for services to an external load balancer. Internally, the swarm lets you specify how to distribute service containers between nodes.

- Secure by default: Each node in the swarm enforces TLS mutual authentication and encryption to secure communications between itself and all other nodes. You have the option to use self-signed root certificates or certificates from a custom root CA.

- Rolling updates: At rollout time you can apply service updates to nodes incrementally. The swarm manager lets you control the delay between service deployment to different sets of nodes. If anything goes wrong, you can roll-back a task to a previous version of the service.



### @TODO Demo time!

#### creating new servers

#### adding nodes to cluster

#### build service

#### deploy service to cluster

#### howto schedule




# Kubernetes

Kubernetes is an open-source platform for automating deployment, scaling, and operations of application containers across clusters of hosts.

It comes from Google.
