# BoulderWelt Crowdedness

## Purpose

Already before Covid-19 all BoulderWelt gyms had an indicator of how crowded the gym is on their website.
Now that they have to limit how many persons are in the gym to ensure distance regulations and so on and so forth this becomes even more important.

I created a simple Prometheus exporter to monitor the current level of _crowdedness_ of all gyms and to see if there are any people waiting in the queue already.

The main purpose of this is to get some insights when it is relatively safe to go to the gym because it's empty and when you definitely should not go to the gym.

## Stack

The whole stack consists of the following components:

* Custom exporter
* Prometheus
* Grafana
* Traefik

The `docker-compose.yml` in this repository is almost the same I am using to run the stack on [www.when2boulder.de](https://www.when2boulder.de).
If you want to run it on your own, feel free to deploy it on your own!

## Building

Right now I haven't had time to setup a CI pipeline, provide a Makefile or any other fancy stuff.
If you want to build the app on your own a simple `go build -o bw-crowdedness ./...` (in case of Linux/Mac) or `go build -o bw-crowdedness.exe ./...` (in case of Windows) will provide you a binary.

Although the recommended way to build/run the app is with Docker.
The `Dockerfile` included in the repository is _self-contained_ and does not need any Go SDK installed on your computer - just Docker.