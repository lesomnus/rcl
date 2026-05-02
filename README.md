# rcl

Go bindings for the ROS 2 Client Library (rcl).

## Setup

Before building any Go project that imports this module, a `rcl-go.pc`
pkg-config file must be present on `PKG_CONFIG_PATH` so that `cgo` can locate
the ROS 2 headers and libraries.

1. Add `gen-pc` as a tool in your project

	```sh
	go get -tool github.com/lesomnus/rcl/cmd/gen-pc
	```

2. Generate `rcl-go.pc`

	Source your ROS 2 environment first (this lets the generator auto-detect the
	active distro), then run:

	```sh
	source /opt/ros/<distro>/setup.bash

	go tool gen-pc -o ~/.local/lib/pkgconfig/rcl-go.pc
	# or redirect to a file with sudo
	go tool gen-pc | sudo tee /usr/local/lib/pkgconfig/rcl-go.pc > /dev/null
	```

	The generator will look for the ROS 2 prefix in the following order:

	1. `AMENT_PREFIX_PATH` environment variable (set by `setup.bash`)
	2. `ROS_DISTRO` environment variable
	3. `/opt/ros/` scan, preferring newer distros

2. Export `PKG_CONFIG_PATH` then build your project

	```sh
	export PKG_CONFIG_PATH="$HOME/.local/lib/pkgconfig:$PKG_CONFIG_PATH"
	CGO_ENABLED=1 go build ./...
	```

### Development

Ignore the above steps if you are developing `rcl-go` itself.
Instead, simply run the script:

```sh
$ ./scripts/config.sh
```

It will generate the `rcl-go.pc` file in the `.pc` directory and everything will be ready for development and testing if you work in the devcontainer.
