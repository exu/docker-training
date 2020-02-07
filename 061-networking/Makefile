##
# Dockerfiles
#
# @file
# @version 0.1

build:
	docker build -t app1 -f Dockerfile.app1 .
	docker build -t app2 -f Dockerfile.app2 .

run:
	docker network create simple-apps || true
	docker run --rm -d -it --network simple-apps --name app1 app1
	docker run --rm -d -it --network simple-apps --name app2 app2

kill:
	docker kill app1 app2

# end
