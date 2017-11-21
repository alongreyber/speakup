start: # Start minikube and dev related deployments
	minikube start --cpus 2	
	minikube addons enable heapster; minikube addons enable ingress
	# Put everything that doesn't run the actual application
	# in the dev namespace. This command creates a local registry we
	# can push docker images to
	kubectl create namespace dev || true
	kubectl --namespace dev apply -f local-registry.yml
	
update: # Update running kubernetes cluster with current code
	# Start local docker registry forwarding container
	docker run -d -e "REGIP=`minikube ip`" -p 30400:5000 chadmoon/socat:latest bash -c "socat TCP4-LISTEN:5000,fork,reuseaddr TCP4:`minikube ip`:30400" || true
	# docker build, docker push everything
	docker build -t 127.0.0.1:30400/web:latest -f web/Dockerfile web
	docker push 127.0.0.1:30400/web:latest
	docker build -t 127.0.0.1:30400/gentle:latest -f gentle/Dockerfile gentle
	docker push 127.0.0.1:30400/gentle:latest
	# Deleting all pods makes them re-pull the latest image
	kubectl delete --all --grace-period=0 --force=true pods
	kubectl apply -f manifests/

install:
	# Docker (Latest stable version)
	sudo apt-get update
	sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
	sudo apt-get update && sudo apt-get install docker-ce
	# Minikube v0.18.0
	curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.18.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/
	# Kubectl (Latest)
	curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
	chmod +x kubectl && sudo mv kubectl /usr/local/bin/
