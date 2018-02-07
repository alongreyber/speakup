start: # Start minikube and dev related deployments
	minikube delete || true
	minikube start --cpus 2	--memory 8192
	# Put everything that doesn't run the actual application
	# in the dev namespace. This command creates a local registry we
	# can push docker images to
	kubectl create namespace dev | true
	kubectl create namespace tiller | true
	helm init --tiller-namespace tiller
	# Wait until tiller starts
	kubectl delete storageclass standard | true
	helm dep update manifests/local-docker-registry
	# Keep retrying because this won't work until the tiller pod is up
	until helm install --tiller-namespace tiller --namespace dev --name dev manifests/local-docker-registry; do sleep 5; done
	# Forwards HTTP to HTTPS as required by docker registry
	# TODO make this not necessary. It will be a lot of work
	docker run -d -e "REGIP=`minikube ip`" -p 30400:5000 chadmoon/socat:latest bash -c "socat TCP4-LISTEN:5000,fork,reuseaddr TCP4:`minikube ip`:30400" | true
	# -------------------------------------
	# Only run once per development session
	# -------------------------------------
	

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
	helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/
	helm repo add another-incubator http://storage.googleapis.com/kubernetes-charts-incubator/
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/
	helm dep update manifests/audio-storage
	helm delete --tiller-namespace tiller $(shell helm --tiller-namespace tiller list | grep audio-storage | cut -f1) | true
	kubectl delete --all pv && kubectl delete --all pvc && kubectl delete --all configmaps && kubectl delete --all statefulsets
	helm install --tiller-namespace tiller manifests/audio-storage
	

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
	# Virtualbox (5.2)
	echo "deb http://download.virtualbox.org/virtualbox/debian xenial contrib" > /etc/apt/sources.list
	wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
	sudo apt update && sudo apt install virtualbox-5.2
	# Minikube v0.23.0 - contains kubernetes server version v1.8.0
	curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.24.1/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
	# Kubernetes client v1.8.7
	curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.8.7/bin/linux/amd64/kubectl
	chmod +x kubectl && sudo mv kubectl /usr/local/bin/
	# Helm v2.8.0 rc1
	curl -o helm.tar.gz https://storage.googleapis.com/kubernetes-helm/helm-v2.8.0-linux-amd64.tar.gz
	tar -zxvf helm.tar.gz
	sudo mv linux-amd64/helm /usr/local/bin/helm
	rm -rf linux-amd64/
	rm helm.tar.gz
	# Restart minikube to apply addon changes
	minikube start
	minikube addons enable heapster; minikube addons enable ingress; minikube addons disable default-storageclass
	minikube delete
