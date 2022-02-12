# ixia-c-one


### Get Started

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
