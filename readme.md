# ixia-c-one


### Get Started

- Run Tests

    ```sh
    cd tests
    # modify hostnames of ixia-c-one (otg) or ceos (dut) if there was a change in .clab.yaml
    vi const.go
    go test -run TestBGPRouteInstall
    ```
