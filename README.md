# SpeakUp
## It's time for your business to SpeakUp 
### Architecture
SpeakUp is built using a distributed containerized system of microservices.
They are organized using Google's Kubernetes system. A local development install includes the entire architecture present on the production system. The process for local development is that you modify your code, push it to a local Docker registry and then the code is rolled out to the cluster. All of this is automated with the Makefile but it is definitely important to know how it works.

### Setting up a development environment
#### Requirements
Tested using Ubuntu 16.04 virtualized in VMWare Workstation. However since everything runs in Docker or Virtualbox you should be able to get it running on basically anything (Except WSL, tried that already).
#### Running the code
There are only three steps:
1. make install
1. make start
1. make update

The first time you run make update it will take a while because it is building the images (specifically the Gentle image takes forever). But after that you should run make update every time you make a code change and it should be near instantaneous (besides the time it takes for the containers to restart).
