# How do Docker plugins work

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

The file has to have the same name as the `--volume-driver` variable in the `docker run` command

!SUB
## Plugin discovery
Once the Docker daemon finds a matching `.sock`, `.spec` or `.json` it will use the `.sock` file or the address as described in the `.spec` or `.json` file to do a HTTP `POST` to on relevant plugin events.

!SUB
## Plugin activation
When a plugin is used for the first time it has to be activated.
Docker will `POST` an empty request to the `/Plugin.Activate` endpoint.

The plugin has to respond with a list of the Docker subsystems the plugin supports:
```
{
    "Implements": ["VolumeDriver"]
}
```
<small>There is currently no way to (manually) deactivate a plugin</small>

!SUB
## Volume plugin details
Apart from activating a Volume plugin also has to respond to `POST`s to the following endpoints:

- VolumeDriver.Create: On creation of a container
- VolumeDriver.Path: To return the path to mount
- VolumeDriver.Mount: On container start, should return the path to mount
- VolumeDriver.Unmount: On container stop
- VolumeDriver.Remove: On `docker rm -v`

!SUB
## Volume plugin details
All responses should include `"Err": null` when there are no errors.

`VolumeDriver.Mount` and `VolumeDriver.Path` should return the path to mount:
<br>`"Mountpoint": "/path/to/directory/on/host"`

For more details see:
<br>[Extend Docker - Volume plugins](https://docs.docker.com/engine/extend/plugins_volume/)

!SUB
## Example
Start the requestlogger Volume plugin
```bash
docker run -ti --rm -v "$PWD":/go/src/myapp -w /go -p 8080:8080 golang:1.5 go run src/myapp/requestlogger.go
```
Start a container using the requestlogger Volume plugin
```bash
docker run -ti -v volumename:/data --volume-driver=requestlogger debian:jessie gliderlabs/alpine
```

!SUB
## Example output
The requests will be logged to `STDOUT` by the requestlogger plugin
```
Request: POST to /VolumeDriver.Create with body: {"Name":"volumename"}
Request: POST to /VolumeDriver.Path with body: {"Name":"volumename"}
Request: POST to /VolumeDriver.Mount with body: {"Name":"volumename"}
Request: POST to /VolumeDriver.Unmount with body: {"Name":"volumename"}
```