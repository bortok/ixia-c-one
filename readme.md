# ixia-c-one


### Get Started

- Run Tests

    ```sh
    # all tests, test configs and helpers are kept inside this directory
    cd tests
    # modify hostnames of ixia-c-one (otg) or ceos (dut) if there was a change in .clab.yaml
    vi const.go
    # modify test contents if needed
    vi bgp_route_install_test.go
    # start the test
    go test -run TestBGPRouteInstall
    ```
