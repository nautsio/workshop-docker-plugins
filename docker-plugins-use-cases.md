# Plugin use-cases

!SUB
## Plugin examples
- Managing networked storage on the host ([NFS](https://github.com/SvenDowideit/docker-volumes-nfs), [Ceph](https://github.com/contiv/volplugin), [GlusterFS](https://github.com/calavera/docker-volume-glusterfs), etc.)
- Moving data volumes between hosts on container migration ([Flocker](https://docs.clusterhq.com/en/1.7.2/introduction/flocker-plugin.html))
- Manage container networking ([Weave](https://github.com/weaveworks/docker-plugin))
- [Custom](https://docs.docker.com/engine/extend/plugins/) storage and networking functionality

!SUB
## Plugin functionality
We are going to implement custom storage functionality in our volume plugin. To keep things simple we'll just manipulate filesystem directories on the host system. It will work just like regular Docker volumes, but we get to determine where the volumes are placed.

!SUB
## Examine and run the plugin
Please open the ```plugin.go``` source file from the part4/ subdirectory in your favourite editor and examine its contents. Do you understand what's going on? If not, discuss! Let's try to run it:
```bash
# Create a container from the 1.5 golang image. If the image is not present it will be downloaded.
# When the container starts it will compile and run the plugin within the container.
#
# Three volumes are bind-mounted into the container:
#
# /run/docker/plugins      : to allow communication between the Docker daemon and the plugin.
# /tmp:/tmp                : used as base path (/tmp/docker/volumes) for new volumes
# Current directory ($PWD) : to build and run the plugin within the container

docker run -it --rm -v /run/docker/plugins:/run/docker/plugins -v /tmp:/tmp -v "$PWD":/go/src/ddd -w /go golang:1.5 go run src/ddd/plugin.go
```

!SUB
## Test the plugin
Let's use the Docker client to actually create containers and use our volume plugin.

```bash
# Run a container, mount a volume to /data and enter the shell
# Examine the log output from the plugin, what do you see?
docker run -it --volume-driver=ddd -v ddd_volume:/data gliderlabs/alpine /bin/sh

# Examine the volume and create a file
ls /data
echo 'Some content for my awesome volume plugin' > /data/readme.txt
# Add additional content in /data if you want
exit

# Run a new container, mount the same volume and enter the shell
docker run -it --volume-driver=ddd -v ddd_volume:/data gliderlabs/alpine /bin/sh
ls /data

# What do you see?
# Play around with the docker volume command to inspect the volume
# and investigate your host to see if the volume is created under
# /tmp/docker/volumes (you might need to enter the boot2docker vm)
```

!SUB
## Final achievement unlocked!
We have managed to run a server that will respond to all volume requests from the Docker daemon. It manipulates the filesystem to create and destroy volumes. This is a simple example but you could extend it with any functionality you require.
