# Docker plugin discovery and activation

!SUB
## What happens under the hood
When you use a Docker plugin
```
docker run -v <volumename>:<mountpoint> --volume-driver=<pluginname> ...
```

!SUB
## Plugin discovery
The Docker daemon looks in order in:
- `/run/docker/plugins` for UNIX domain sockets
- `/etc/docker/plugins` for `.spec` or `.json` files
- `/usr/lib/docker/plugins` for `.spec` or `.json` files

The file has to have the same name as `<pluginname>` in the `docker run` command from the previous slide.

!SUB
## Plugin discovery
Once the Docker daemon finds a matching `.sock`, `.spec` or `.json` it will use the `.sock` file or the address as described in the `.spec` or `.json` file to do a HTTP `POST` on relevant plugin events.

!SUB
## Example spec file
This JSON structure describes a Docker plugin:
```json
{
  "Name": "plugin-example",
  "Addr": "https://example.com/docker/plugin",
  "TLSConfig": {
    "InsecureSkipVerify": false,
    "CAFile": "/usr/shared/docker/certs/example-ca.pem",
    "CertFile": "/usr/shared/docker/certs/example-cert.pem",
    "KeyFile": "/usr/shared/docker/certs/example-key.pem",
  }
}
```

!SUB
## Plugin activation
When a plugin is used for the first time it has to be activated.
Docker will `POST` an empty request to the `/Plugin.Activate` endpoint.

The plugin has to respond with a list of the Docker subsystems the plugin supports:
```json
{
    "Implements": ["VolumeDriver"]
}
```
<small>There is currently no way to (manually) deactivate a plugin</small>

!SUB
## Plugin Activation
From the previous slides we have learned that the plugin needs to be activated by the Docker daemon.
Let's extend our server with a new activation route.

!SUB
## Examine and run the plugin
Please open the ```plugin.go``` source file from the part2/ subdirectory in your favourite editor and examine its contents. Do you understand what's going on? If not, discuss! Let's try to run it:
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
# In a seperate terminal create a container from the Nauts toolbox-networking and run a shell.
# This image is based on Alpine and additionally contains the curl and dig commands.
#
# The /run/docker/plugins directory is bind-mounted into the container
# to allow access to the plugin socket.
docker run -it --rm -v /run/docker/plugins:/run/docker/plugins cargonauts/toolbox-networking /bin/sh

# Run a HTTP request with an empty JSON body against the socket on the /Plugin.Activate endpoint
curl -s -XPOST --unix-socket /run/docker/plugins/ddd.sock -d "{}" http:/Plugin.Activate

# This should report "{"Implements":["VolumeDriver"]}" in the response as an indication
# that the plugin can be activated by the Docker daemon.
```

!SUB
## Achievement unlocked!
We have managed to run a server that listens for incoming HTTP requests and will respond to an activation call from the Docker daemon!
