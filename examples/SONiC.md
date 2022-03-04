# Test examples

## SONiC

1. Create Topology

```sh
sudo containerlab deploy --topo ixia-c-one-sonic.clab.yaml
```

2. Configure SONiC node

````
docker exec -it clab-ixia-c-sonic bash

sysctl net.ipv6.conf.Ethernet0.disable_ipv6=0
config interface ip add Ethernet0 1.1.1.3/24
config interface ip add Ethernet0 0:1:1:1::3/64
config interface startup Ethernet0
ip link set eth1 up

sysctl net.ipv6.conf.Ethernet4.disable_ipv6=0
config interface ip add Ethernet4 2.2.2.3/24
config interface ip add Ethernet4 0:2:2:2::3/64
config interface startup Ethernet4
ip link set eth2 up

config loopback add Loopback0
sysctl net.ipv6.conf.Loopback0.disable_ipv6=0
config interface ip add Loopback0 3.3.3.3/32
config interface ip add Loopback0 0:3:3:3::3/64
config interface startup Loopback0

cp /etc/frr/daemons /etc/frr/daemons.orig
cat /etc/frr/daemons.orig | sed "s/^bgpd=no$/bgpd=yes/" > /etc/frr/daemons

service frr restart
service frr status | grep bgpd

vtysh
configure
frr defaults traditional
ipv6 forwarding
router bgp 3333
  bgp router-id 3.3.3.3
  neighbor 1.1.1.1 remote-as 1111
  neighbor 2.2.2.2 remote-as 2222
  neighbor 0:1:1:1::1 remote-as 1111
  neighbor 0:2:2:2::2 remote-as 2222
  no bgp ebgp-requires-policy
  no bgp network import-check
  address-family ipv4 unicast
    neighbor 1.1.1.1 soft-reconfiguration inbound
    neighbor 1.1.1.1 activate
    neighbor 2.2.2.2 soft-reconfiguration inbound
    neighbor 2.2.2.2 activate
    exit-address-family
  address-family ipv6 unicast
    neighbor 0:1:1:1::1 activate
    neighbor 0:1:1:1::1 soft-reconfiguration inbound
    neighbor 0:2:2:2::2 activate
    neighbor 0:2:2:2::2 soft-reconfiguration inbound
    exit-address-family
  exit
exit
sh ip bgp summary
exit
exit
````


5. Run Tests

    ```sh
    # all tests, test configs and helpers are kept inside this directory
    cd tests
    # modify hostnames of ixia-c-one (otg) or ceos (dut) if there was a change in .clab.yaml
    vi const.go
    
    # Modify test contents of L3 forwarding test with DUT acting as BGP router if needed and note the
    # name of Test* function
    vi bgp_route_install_test.go
    # Run the test using the name noted above
    go test -run TestSONiCIPv4BGPRouteInstall
    go test -run TestSONiCIPv6BGPRouteInstall
    ```



6. Destroy Topology

```sh
sudo containerlab destroy --topo ixia-c-one-sonic.clab.yaml --cleanup
```
