start: # Start minikube and dev related deployments
	minikube start --cpus 2	--memory 4096
	minikube addons enable heapster; minikube addons enable ingress
	# Put everything that doesn't run the actual application
	# in the dev namespace. This command creates a local registry we
	# can push docker images to
	kubectl create namespace dev | true
	kubectl create namespace tiller | true
	helm init --tiller-namespace tiller
	# Wait until tiller starts
	until kubectl get deployments --namespace tiller -o=jsonpath='{range .items[*]}{.status.availableReplicas}' | grep 1; do sleep 1; done
	# TODO why doesn't this wait long enough on it's own?
	sleep 10
	helm install --namespace dev --tiller-namespace tiller --name dev manifests/local-docker-registry
	# Only run once per development session
	
update: # Update running kubernetes cluster with current code
	./sass/update_css.sh
	# docker build, docker push everything
	# web, gentle, audio-transcoder, kafka-connector
	docker build -t 127.0.0.1:30400/web:latest -f web/Dockerfile web
	docker push 127.0.0.1:30400/web:latest
	docker build -t 127.0.0.1:30400/gentle:latest -f gentle/Dockerfile gentle
	docker push 127.0.0.1:30400/gentle:latest
	docker build -t 127.0.0.1:30400/audio-transcoder:latest -f audio-transcoder/Dockerfile audio-transcoder
	docker push 127.0.0.1:30400/audio-transcoder:latest
	docker build -t 127.0.0.1:30400/kafka-postgres-connector:latest -f kafka-postgres-connector/Dockerfile kafka-postgres-connector
	docker push 127.0.0.1:30400/kafka-postgres-connector:latest

install: # Install dependencies. No checking for if they are already installed. Must be run as root.
	sudo apt update
	sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
	# Docker (latest stable version)
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo add-apt-repository \
	    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
	    $(shell lsb_release -cs) \
	    stable"
	sudo apt update && sudo apt install -y docker-ce
	# Enable docker for non-root
	sudo groupadd docker | true
	sudo gpasswd -a ${USER} docker
	newgrp docker
	# Virtualbox (5.2)
	echo "deb http://download.virtualbox.org/virtualbox/debian xenial contrib" > /etc/apt/sources.list
	wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
	sudo apt update && sudo apt install virtualbox-5.2
	# Minikube v0.18.0 - contains kubernetes server version v1.6
	curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.24.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
	# Kubernetes client v1.6
	curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.6.13/bin/linux/amd64/kubectl
	chmod +x kubectl && sudo mv kubectl /usr/local/bin/
	# Helm
	curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash
