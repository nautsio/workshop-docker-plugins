# Docker plugin communication

!SUB
## Volume plugin details
Apart from activating, a volume plugin also has to respond to `POST` requests to the following endpoints:

- VolumeDriver.Create: On creation of a container
- VolumeDriver.Path: To return the mount path
- VolumeDriver.Mount: On container start, return the mount path
- VolumeDriver.Unmount: On container stop
- VolumeDriver.Remove: On volume removal through the Docker client

!SUB
## Volume plugin details
All responses should include `"Err": null` in their response when there are no errors.

The VolumeDriver.Mount and VolumeDriver.Path routes should also return the mount path:
```json
{
    "Mountpoint": "/path/to/directory/on/host",
    "Err": null
}
```

!SUB
## Plugin routes
From the previous slides we have learned that the plugin needs additional routes to handle all the Docker daemon volume requests.
Let's extend our server with these new routes.

!SUB
## Examine and run the plugin
Please open the ```plugin.go``` source file from the part3/ subdirectory in your favourite editor and examine its contents. Do you understand what's going on? If not, discuss! Let's try to run it:
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
Let's use the Docker client volume subcommand to interact with our plugin. Examine the output from our plugin after each command. What do you see?
```bash
# Create volume with our custom driver
docker volume create --driver=ddd --name=ddd_plugin_test

# Inspect the volume properties
# Notice that inspect output does not list a mountpath yet.
docker volume inspect ddd_plugin_test

# Remove the volume
docker volume rm ddd_plugin_test
```

!SUB
## Achievement unlocked!
We have managed to run a server that listens for incoming HTTP requests and will respond to all volume requests from the Docker daemon!
