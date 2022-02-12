# ixia-c-one

`ixia-c-one` is a re-packaged (as a single-container) flavor of multi-container application [ixia-c](https://github.com/open-traffic-generator/ixia-c).  
This repository hosts bare minimum artifacts (configurations and tests) to get started with [containerlab](https://containerlab.srlinux.dev/) and `ixia-c-one`.

### Prerequisites

- x86-64 Ubuntu 20.04 Server
- At least 2 CPU cores, 4GB RAM and 64GB HDD
- Docker (https://docs.docker.com/engine/install/ubuntu/)
- Go 1.17+ (https://go.dev/doc/install)
- curl

### Get Started

- Get containerlab with added support for ixia-c-one.

    ```sh
    curl -kLO https://github.com/open-traffic-generator/ixia-c-one/releases/download/v0.0.1-2610/containerlab
    chmod +x containerlab
    ```

- Create Topology

    ```sh
    sudo ./containerlab deploy --topo ixia-c-one-ceos.clab.yaml
    ```

    > If this step fails, most probably you do not have the ceos docker image. 
    > Please obtain the image from https://www.arista.com/en/support/software-download and re-tag it as specified in .yaml.

- Run Tests

    ```sh
    # all tests, test configs and helpers are kept inside this directory
    cd tests
    # modify hostnames of ixia-c-one (otg) or ceos (dut) if there was a change in .clab.yaml
    vi const.go
    # modify test contents if needed and note the name of Test* function
    vi bgp_route_install_test.go
    # run the test using the name noted above
    go test -run TestBGPRouteInstall
    ```

- Destroy Topology

    ```sh
    sudo ./containerlab destroy --topo ixia-c-one-ceos.clab.yaml
    ```
