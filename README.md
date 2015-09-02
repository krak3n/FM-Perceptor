# SOON\_ FM | Perceptor

<img src="perceptor.jpg" width="80" height="114" align="right" />

Perceptor is the new SFM\_ events system for websocket communication between the Player and API.
This allows the player to connect to the events system without the need for a direct connection
to Redis, meaning the player no longer needs to connect to the VPN to get events.

The ultilate goal will be to remove Redis as the main events hub, moving to REST pushes to various
evented services on demand.

<div align="center">
    <img src="infrastructor.png" />
</div>

## Running

Perceptor can be run via Docker. Simply download the image and run:

```
docker pull registry.soon.build/fm/perceptor
docker run --rm -it registry.soon.build/fm/perceptor
```

### Configuration

Perceptor can be configured via a config file located in any of the following locations:

* `/etc/perceptor/perceptor.yml`
* `$HOME/.perceptor/perceptor.yml`
* `$PWD/.perceptor/perceptor.yml`

The config file atleast needs to set client secrets for HMAC request verification:

``` yaml
clients:
  soundwave: foo
  shockwave: bar
```

In addition Redis connection settings and the port `perceptor` will run on can also be configured:

``` yaml
# Port to run perceptor on
port: 9000
# The log verbosity level (debug, info, wan, error)
log_leveL: debug
# Redis Host Address
redis_host: 127.0.0.1
# Redis Port
redis_port: 6379
# Client Secrets - For HMAC Verification
clients:
  soundwave: foo
  shockwave: bar
```

You can share this file with the docker container, for example:

```
docker run --rm -it -v /path/to/perceptor.yml:/etc/perceptor/perceptor.yml registry.soon.build/fm/perceptor
```

#### Environment Variables

In addition the following environment variables can override configuration:

* `PERCEPTOR_PORT` - The port `perceptor` runs on
* `PERCEPTOR_LOG_LEVEL` - The verbosity of the logging (`debug`, `info`, `warn`, `error`)
* `PERCEPTOR_REDIS_HOST` - Redis Host Address (`localhost`)
* `PERCEPTOR_REDIS_PORT` - Redis Port (`6379`)

## Development

This package uses [Glide](https://github.com/Masterminds/glide) for vendoring, please follow the
install instructions on the [Glide](https://github.com/Masterminds/glide) repository. Ensure you
have `Go 1.5` and `Glide 0.5.0` installed.

### Create a Workspace

Firsrt create a workspace, this will form part of your `$GOPATH`.

```
mkdir -p ~/Development/Perceptor/src/github.com/thisissoon/FM-Perceptor
```

Now set your `$GOPATH`:

```
export GOPATH=~/Development/Perceptor
export GO15VENDOREXPERIMENT=1
```

We have also set `GO15VENDOREXPERIMENT=1` to make use of the `/vendor` dirctory support in `Go 1.5`.

### Clone the Project

Now you can clone the project into the workspace.

```
cd ~/Development/Perceptor/src/github.com/thisissoon/FM-Perceptor
git clone git@github.com:thisissoon/FM-Perceptor.git .
```

### Get Dependencies

Now we can get the dependencies using `glide`.

```
glide up
```

### Build

Now the project can be built:

```
go build $(glide nv)
```
