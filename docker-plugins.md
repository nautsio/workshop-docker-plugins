# Docker plugins

!SUB
## Docker plugins
As of version [1.7.0](https://blog.docker.com/2015/06/announcing-docker-1-7-multi-host-networking-plugins-and-orchestration-updates/) Docker has experimental support for plugins for [networking](https://docs.docker.com/engine/extend/plugins_network/) and [volumes](https://docs.docker.com/extend/plugins_volume/)(storage).
<br><br>In Docker [1.8.0](https://blog.docker.com/2015/08/docker-1-8-content-trust-toolbox-registry-orchestration/) the volumes plugin was promoted to the stable release and Docker [1.9.0](https://blog.docker.com/2015/11/docker-1-9-production-ready-swarm-multi-host-networking/) saw the network plugin promoted to the stable release.

!SUB
## What are Docker plugins?
Docker plugins are an easy way for third parties to extend Docker's functionality.
<br><br>
Docker plugins are specific to a Docker subsystem. At the moment there is support for plugins in the [network subsystem ](https://github.com/docker/libnetwork/blob/master/docs/remote.md) and the [volume subsystem](https://docs.docker.com/engine/extend/plugins_volume/)

!SUB
## How do Docker plugins work?
Docker plugins are out-of-process programs that expose a webhook-like functionality which the Docker daemon uses to send HTTP POST requests so that the plugin can act on Docker events

!SUB
## How do Docker plugins work?
A docker plugin has to expose either a UNIX domain socket or a HTTP endpoint on the (remote) Docker host so that the Docker daemon can talk to it.

!SUB
## How to use a Docker plugin?
For volume plugins simply use:
```
docker run -v <volumename>:<mountpoint> --volume-driver=<pluginname> ...
```

For network plugins simply create the network first and then use it
```
docker network create -d <pluginname> <networkname>
docker run --net=<networkname> ...
```
<small>(this will only work on Docker 1.9.0 and later)</small>
