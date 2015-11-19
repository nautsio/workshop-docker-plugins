# Creating a Docker plugin

!SUB
## Our first Docker plugin
Let's start iteratively building our own Docker volume plugin.

Based on what we've learned so far we know that we need an HTTP server that listens for incoming requests from the Docker daemon

!SUB
## Examine and run the plugin
Please open the ```plugin.go``` files from the part1/ subdirectory in your favourite editor and examine its contents. Notice that the plugin communicates through a socket in ```/run/docker/plugins```. Discuss the code if needed. <br><br>Let's try to run it:
```bash
# Create a container from the 1.5 golang image. If the image is not present it will be downloaded.
# When the container starts it will compile and run the plugin within the container.
#
# Three volumes are bind-mounted into the container:
#
# /run/docker/plugins      : to allow communication between the Docker daemon and the plugin.
# Current directory ($PWD) : to build and run the plugin within the container

docker run -it --rm -v /run/docker/plugins:/run/docker/plugins -v "$PWD":/go/src/ddd -w /go golang:1.5 go run src/ddd/plugin.go
```

!SUB
## Test the plugin
Let's test our plugin:
```bash
# In a seperate terminal window create a container from the Nauts toolbox-networking and run a shell.
# This image is based on Alpine and additionally contains the curl and dig commands.
#
# The /run/docker/plugins directory is bind-mounted into the container
# to allow access to the plugin socket.
docker run -it --rm -v /run/docker/plugins:/run/docker/plugins cargonauts/toolbox-networking /bin/sh

# Run a HTTP request with an empty JSON body against the socket on the / endpoint
curl -s --unix-socket /run/docker/plugins/ddd.sock -d "{}" http:/

# This should report "Dutch Docker Day" in the response as an indication
# that the HTTP server works!
```

!SUB
## Achievement unlocked!
We have managed to run a server that listens for incoming HTTP requests! Now we can start to actually make it listen for incoming Docker daemon events.
