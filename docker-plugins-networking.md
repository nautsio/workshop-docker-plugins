# Networking plugins

!SUB
## Docker Network Plugins
So far we have looked in detail at volume plugins. In the last part of this workshop we will take a closer look at the low-level implementation of network plugins and see how they are different (or not) from volume plugins.

!SUB
## Network plugin architecture
Architecturally Docker network plugins are identical to volume plugins. They are run as separate out-of-band processes that run a HTTP server and respond to requests from the Docker daemon.

!SUB
## Network plugin activation
Like the volume plugin, a network plugin needs to activate itself with the Docker daemon before it can be used. The Docker daemon will issue the same '/Plugin.Activate' HTTP POST request to the plugin and it should respond with:
```json
{
    "Implements": ["NetworkDriver"]
}
```
It will also receive an additonal HTTP POST request to '/NetworkDriver.GetCapabilities'. This time the response should be:
```json
{
    "Scope": "local"
}
```
Scope can be either local or global.

!SUB
## Network plugin routes
Apart from activating and announcing capabilities, a network plugin also has to respond to `POST` requests to the following endpoints:

- NetworkDriver.CreateNetwork: to create a new network
- NetworkDriver.DeleteNetwork: to delete a network
- NetworkDriver.CreateEndpoint: to return an endpoint
- NetworkDriver.EndpointOperInfo: to get operational information
- NetworkDriver.DeleteEndpoint: to delete an endpoint
- NetworkDriver.Join: to return the interface, gateway and routes
- NetworkDriver.Leave: to remove and endpoint from the sandbox
- NetworkDriver.DiscoverNew: to receive discovery notifications
- NetworkDriver.DiscoverDelete: to receive the DiscoverDelete notification

!SUB
## Network plugin use-cases
So what kind of networks could we create with a network plugin?

- bridged: this is default Docker mode (docker0 bridge)
- [vxlan](https://en.wikipedia.org/wiki/Virtual_Extensible_LAN): used by the Docker [overlay driver](https://docs.docker.com/engine/userguide/networking/get-started-overlay/) and Weave [fast-datapath](http://blog.weave.works/2015/06/12/weave-fast-datapath/) to allow multi-host networking
- [macvlan](http://www.pocketnix.org/posts/Linux%20Networking:%20MAC%20VLANs%20and%20Virtual%20Ethernets):  A MAC VLAN takes a single network interface and creates multiple virtual ones with different MAC addresses ([example](https://github.com/gopher-net/macvlan-docker-plugin))
- [ipvlan](https://coreos.com/rkt/docs/latest/networking.html): ipvlan creates virtual copies of interfaces like macvlan but does not assign a new MAC address to the copied interface ([example](https://github.com/gopher-net/ipvlan-docker-plugin))
- custom: you can create any custom networking configuration for your containers that you like
