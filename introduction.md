# Introduction

!SUB
## Menu for today
- What are Docker plugins and how do they work
- Which Docker plugins are available
- Build your own Docker plugin!

!SUB
## Hands-on guidelines
- Form groups of 2-3 people
- Help each other when you can
- Use the slides as guidance


!SLIDE
# Setup

!SUB
## Prerequisites
[Docker](https://docs.docker.com/installation) <small>(v1.8.0+)</small>

[Go Docker image](https://hub.docker.com/r/library/golang) <small>(golang:1.5+)</small>

!SUB
## Get the files

[github.com/nautsio/workshop-docker-plugins](https://github.com/nautsio/workshop-docker-plugins)
```
$ git clone https://github.com/nautsio/workshop-docker-plugins.git
$ cd workshop-docker-plugins
```
Or download the files directly
<br>[zip](https://github.com/nautsio/workshop-docker-plugins/archive/master.zip) or [tar.gz](https://github.com/nautsio/workshop-docker-plugins/archive/master.tar.gz)


!SLIDE
# Docker plugins

!SUB
## Docker plugins
As of version [1.7.0](https://blog.docker.com/2015/06/announcing-docker-1-7-multi-host-networking-plugins-and-orchestration-updates/) Docker has experimental support for plugins for [networking](https://github.com/docker/libnetwork/blob/master/docs/remote.md) and [volumes](https://docs.docker.com/extend/plugins_volume/)(storage).
As of Docker [1.8.0](https://blog.docker.com/2015/08/docker-1-8-content-trust-toolbox-registry-orchestration/) the volumes plugin has been promoted to the stable release. The network plugin will be promoted to the stable release with Docker 1.9.0

!SUB
## What are Docker plugins?
Docker plugins are an easy way for third parties to extend Docker´s functionality.

!SUB
## What are Docker plugins?
Docker plugins are specific to a Docker subsystem.
At the moment there is support for plugins in the [network subsystem ](https://github.com/docker/libnetwork/) and the volume subsystem

!SUB
## How do Docker plugins work?
Docker plugins are out-of-process programs that expose a webhook-like functionality which the Docker daemon uses to `POST` HTTP requests to so the plugin can act on Docker events

<small>Currently running plugins outside containers is recommended, because plugins should be started before and stopped after the Docker daemon</small>

!SUB
## How do Docker plugins work?
A docker plugin has to either expose a UNIX domain socket on the Docker host or a HTTP endpoint on the Docker host or on a remote host

!SUB
## How to use a Docker plugin?
For volume plugins simply use:
```
docker run --volume-driver=<pluginname> ...
```

For network plugins simply create the network first and then use it
```
docker network create -d <pluginname> <networkname>
docker run --net=<networkname> ...
```
<small>(this will only work on Docker 1.9.0 and later)</small>

!SUB
## Discovery of Docker plugins
Plugins are discovered by name and a simple check for a file with the same name in specific directories
Docker always searches for a UNIX domain socket first in `/run/docker/plugins`. If the socket is not found it will continue to check for `.spec` or `.json` files in `/etc/docker/plugins` and `/usr/lib/docker/plugins`.

Fore more details see [Docker docs - Plugin discovery](https://docs.docker.com/extend/plugin_api/#plugin-discovery)

!SUB
## Caveats
Plugins have to be started before the Docker daemon or there will be a slight delay before the docker daemon is started because the Docker daemon tries to connect to the plugin while starting. Docker will output the following in `docker.log`
```
Unable to connect to plugin: 127.0.0.1:8080, retrying in 1s
Unable to connect to plugin: 127.0.0.1:8080, retrying in 2s
Unable to connect to plugin: 127.0.0.1:8080, retrying in 4s
Unable to connect to plugin: 127.0.0.1:8080, retrying in 8s
```

!SUB
> Attempts to call a method on a plugin are retried with an exponential backoff for up to 30 seconds. This may help when packaging plugins as containers, since it gives plugin containers a chance to start up before failing any user containers which depend on them.

[Docker plugin API docs](https://github.com/docker/docker/blob/master/docs/extend/plugin_api.md#plugin-retries)
