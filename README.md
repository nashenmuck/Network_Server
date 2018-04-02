# Network - Server Component
### What is this?
Network is an implementation of a simple federated social media network with an open API. This allows for anyone to create their own client, and host a server which can communicate with other Network servers.
### Building
Network is built in Go, and can be built with:
```
go get github.com/nashenmuck/network_server
```
A Dockerfile is also supplied and can be built in the root project directory with:
```
docker build -t nashenmuck/network .
```
If you are about to deploy to Minikube, and wish to test your build on there, you may wish to run `eval $(minikube docker-env)` before building the Docker image.
### Deploying
Network is Kubernetes native, and comes with a Helm chart for deployment. It has Postgres as a dependency, which will be installed along with the Network server. It can be installed with:
```
helm install -n network chart/network
```
For values that can be set, please see the `values.yaml` file. Note that values for the Postgres dependency can also be set by doing `--set postgresql.value=foo`. If you're deploying to ARM, you'll need to use at least Postgres version 10.3. You can do this with `--set postgresql.imageTag=10.3` along with the rest of your command to install the Network Helm chart. You'll also need to `--set image.tag=armhf` for ARM deployments. Network Docker images are currently tagged with their commit hashes for each build, with the `latest` tag being on the latest commit from the `master` branch.