# What's Docker 

## Containers intro (container vs VMs)
## Docker service 
## Docker containers repository 
## Docker images as layers

# Docker examples

## Simple first run [Dockerfile](001-simple-dockerfile/Dockerfile)

- Building with `docker build -t TAG_NAME`.
- Running  with `docker run TAG_NAME


## Entrypoints [Dockerfile](002-entrypoint/Dockerfile)

Default entrypoing in docker is `/bin/sh -c` which simply runs command passed to `CMD` instruction


## Inserting editor inside docker 

You can run almost all applications / services in Docker it's encapsulated OS.

## Docker images (ubuntu vs alpine vs debian vs scratch)

You need be careful in your choice of base docker image they can be huge, and you for sure don't want to pass so big images through your network.


## Docker containers are immutable 

If you change something in your docker images it'll simply lost after docker container 
will be reloaded. 

## Attaching volumes

## Creating Web app 
### Golang based (app server)
### PHP based application

## Exposing application






# Docker training basics for dev worklows

We will be talking today about web based applications:
Which of you have made some web applications?

# Introduction

There are so many possibilities to share your software to the world
The most popular Cloud providers will create for You virtal machine
images.

For example when you're using Amazon / Google Cloud Engine / Digital
Ocean you can easily request for creating new node, engine will spawn
new Virtual Machine with OS of your choice and here you are you have
separate Operating System connected to the Internet. These Machines
are spawned on top of powerful servers with hardware virtualisation.

Where is here place for Docker ? WTF is Docker? Maybe simple example:

After create of such VM you can put your app and start your business
next you out another app, and another and another.

Let's name them app1 app2 app3, app1 have some dependency e.g. ssh-server
in version 1.1 but app2 want ssh-server in version 2.0 and SSL version 2.0 because of many
other dependencies there is no way for your OS to handle such deps tree
app3 want ssh-server version 3.0 and ssl version 1.0. Things are getting
little messy (and it's only 2 deps with 3 different versions)

Imagine that you can put app1 isolate it in
her own environment giv her ssh-server in
version 1.1 and SSL in 2.0.

You can do it of course with VMs but, VMs are resource heavy,
big, and sharing them by network is hard.

And here comes Docker.


# Docker basics


1.  Docker basics

-   VM vs containers

[![img](res/docker-training-vm.png)](res/docker-training-vm.png)

[![img](res/docker-training-containers.png)](res/docker-training-containers.png)

2.  Structure: Images vs containers

- Image: like class in programming
- Container: Instance

We're storing, pushing, changing classes, but not runned instances
(we need to rebuild, in most programming langs :) )

3. File system

File system in docker containers is temporary by default (exception here are data volumes. When you stop container and start again all data will be lost
same will happen with new instance of class in OOP. Persistance here is made
on building process.

3.  Containers Repository

There is huge containers repository with official and unofficial images,
here you can find many of prebuild container images with different linux
distributions, web servers, applications, databases, production ready
systems and many many more.

<https://hub.docker.com/>

now let's run some containers

# Running existing containers

But how containers work? How you need to handle them. Let's start with really simple
example:

run ubuntu with sth..

```
docker run ubuntu cat /etc/passwd
docker run ubuntu apt-get
docker run -it ubuntu /bin/bash
```

# Creating First container with CMD

-   Shell app

Add lines to `Dockerfile`
```
FROM ubuntu:latest
CMD date
```
next you'll need to run

```
docker build -t cmd .
docker run cmd
```

When we pass additional parameters We override CMD section.

```
docker run cmd ls -la
```



# Default program to run ENTRYPOINT

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

# Getting our container IP address

Ok we've run our web app, It's working but how to show it in browser? We'll need
our container IP address:

```
docker inspect

docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

```

result:

```
172.17.0.4
```

http://172.17.0.4:8080/allaallalal


# Data volumes

I've said on beggining that docker has temporary file system by default.
But what when you want to store some informations after container will
be removed?


## Mounting local path inside container for devlopment

We can run our container with local volume mounted inside our container

## Shared-storage Volumes

We can also create docker based volumes which are responsible for
saving data

To do it we need to create data volume:

```
docker volume create --name our-data
docker run -v our-data:/vol -it ubuntu /bin/bash
```



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
