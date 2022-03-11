/* Test BGP Route Installation
Topology:
IXIA (0:40:40:40::0/64) -----> SONiC ------> IXIA (0:50:50:50::0/64)
Flows:
- permit v4: 0:40:40:40::1 -> 0:50:50:50::1+
- deny v4: 0:40:40:40::1 -> 0:60:60:60::1+
*/
package tests

import (
	"testing"

	"github.com/open-traffic-generator/ixia-c-one/tests/helpers"
	"github.com/open-traffic-generator/snappi/gosnappi"
)

func TestSONiCIPv6BGPRouteInstall(t *testing.T) {
	client, err := helpers.NewClient(otgHttpLocation)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()
	defer client.StopProtocol()

	config, expected := bgpRouteInstallConfigSONiCIPv6(client)

	if err := client.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	if err := client.StartProtocol(); err != nil {
		t.Fatal(err)
	}

	helpers.WaitFor(t, func() (bool, error) { return client.AllBgp6SessionUp(expected) }, nil)

	ethernetNames := []string{"dutPort1.eth", "dutPort2.eth"}
	if _, err := client.GetIPv6NeighborsStates(ethernetNames); err != nil {
		t.Fatal(err)
	}

	if err := client.StartTransmit(nil); err != nil {
		t.Fatal(err)
	}

	helpers.WaitFor(t, func() (bool, error) { return client.FlowMetricsOk(expected) }, nil)
}

func bgpRouteInstallConfigSONiCIPv6(client *helpers.ApiClient) (gosnappi.Config, helpers.ExpectedState) {
	config := client.Api().NewConfig()

	port1 := config.Ports().Add().SetName("ixia-c-port1").SetLocation(otgPort1Location)
	port2 := config.Ports().Add().SetName("ixia-c-port2").SetLocation(otgPort2Location)

	dutPort1 := config.Devices().Add().SetName("dutPort1")
	dutPort1Eth := dutPort1.Ethernets().Add().
		SetName("dutPort1.eth").
		SetPortName(port1.Name()).
		SetMac("00:00:01:01:01:01")
	dutPort1Ipv6 := dutPort1Eth.Ipv6Addresses().Add().
		SetName("dutPort1.ipv6").
		SetAddress("0:1:1:1::1").
		SetGateway("0:1:1:1::3")
	dutPort2 := config.Devices().Add().SetName("dutPort2")
	dutPort2Eth := dutPort2.Ethernets().Add().
		SetName("dutPort2.eth").
		SetPortName(port2.Name()).
		SetMac("00:00:02:01:01:01")
	dutPort2Ipv6 := dutPort2Eth.Ipv6Addresses().Add().
		SetName("dutPort2.ipv6").
		SetAddress("0:2:2:2::2").
		SetGateway("0:2:2:2::3")

	dutPort1Bgp := dutPort1.Bgp().
		SetRouterId("1.1.1.1")
	dutPort1Bgp6Peer := dutPort1Bgp.Ipv6Interfaces().Add().
		SetIpv6Name(dutPort1Ipv6.Name()).
		Peers().Add().
		SetName("dutPort1.bgp6.peer").
		SetPeerAddress(dutPort1Ipv6.Gateway()).
		SetAsNumber(1111).
		SetAsType(gosnappi.BgpV6PeerAsType.EBGP)

	dutPort1Bgp6PeerRoutes := dutPort1Bgp6Peer.V6Routes().Add().
		SetName("dutPort1.bgp4.peer.rr6").
		SetNextHopIpv6Address(dutPort1Ipv6.Address()).
		SetNextHopAddressType(gosnappi.BgpV6RouteRangeNextHopAddressType.IPV6).
		SetNextHopMode(gosnappi.BgpV6RouteRangeNextHopMode.MANUAL)
	dutPort1Bgp6PeerRoutes.Addresses().Add().
		SetAddress("0:40:40:40::0").
		SetPrefix(64).
		SetCount(5).
		SetStep(2)

	dutPort2Bgp := dutPort2.Bgp().
		SetRouterId("2.2.2.2")
	dutPort2Bgp6Peer := dutPort2Bgp.Ipv6Interfaces().Add().
		SetIpv6Name(dutPort2Ipv6.Name()).
		Peers().Add().
		SetName("dutPort2.bgp6.peer").
		SetPeerAddress(dutPort2Ipv6.Gateway()).
		SetAsNumber(2222).
		SetAsType(gosnappi.BgpV6PeerAsType.EBGP)

	dutPort2Bgp6PeerRoutes := dutPort2Bgp6Peer.V6Routes().Add().
		SetName("dutPort2.bgp4.peer.rr6").
		SetNextHopIpv6Address(dutPort2Ipv6.Address()).
		SetNextHopAddressType(gosnappi.BgpV6RouteRangeNextHopAddressType.IPV6).
		SetNextHopMode(gosnappi.BgpV6RouteRangeNextHopMode.MANUAL)
	dutPort2Bgp6PeerRoutes.Addresses().Add().
		SetAddress("0:50:50:50::0").
		SetPrefix(64).
		SetCount(5).
		SetStep(2)

	// OTG traffic configuration
	f2 := config.Flows().Add().SetName("p1.v6.p2.permit")
	f2.Metrics().SetEnable(true)
	f2.TxRx().Device().
		SetTxNames([]string{dutPort1Bgp6PeerRoutes.Name()}).
		SetRxNames([]string{dutPort2Bgp6PeerRoutes.Name()})
	f2.Size().SetFixed(512)
	f2.Rate().SetPps(500)
	f2.Duration().FixedPackets().SetPackets(1000)
	e2 := f2.Packet().Add().Ethernet()
	e2.Src().SetValue(dutPort1Eth.Mac())
	//e2.Dst().SetValue("02:42:ac:14:14:02")
	e2.Dst().SetValue("00:00:00:00:00:00")
	v6 := f2.Packet().Add().Ipv6()
	v6.Src().SetValue("0:40:40:40::1")
	v6.Dst().Increment().SetStart("0:50:50:50::1").SetStep("::1").SetCount(5)

	f2d := config.Flows().Add().SetName("p1.v6.p2.deny")
	f2d.Metrics().SetEnable(true)
	f2d.TxRx().Device().
		SetTxNames([]string{dutPort1Bgp6PeerRoutes.Name()}).
		SetRxNames([]string{dutPort2Bgp6PeerRoutes.Name()})
	f2d.Size().SetFixed(512)
	f2d.Rate().SetPps(500)
	f2d.Duration().FixedPackets().SetPackets(1000)
	e2d := f2d.Packet().Add().Ethernet()
	e2d.Src().SetValue(dutPort1Eth.Mac())
	//	e2d.Dst().SetValue("02:42:ac:14:14:02")
	e2d.Dst().SetValue("00:00:00:00:00:00")
	v6d := f2d.Packet().Add().Ipv6()
	v6d.Src().SetValue("0:40:40:40::1")
	v6d.Dst().Increment().SetStart("0:60:60:60::1").SetStep("::1").SetCount(5)

	expected := helpers.ExpectedState{
		Bgp6: map[string]helpers.ExpectedBgpMetrics{
			dutPort1Bgp6Peer.Name(): {Advertised: 5, Received: 10},
			dutPort2Bgp6Peer.Name(): {Advertised: 5, Received: 10},
		},
		Flow: map[string]helpers.ExpectedFlowMetrics{
			f2.Name():  {FramesRx: 1000, FramesRxRate: 0},
			f2d.Name(): {FramesRx: 0, FramesRxRate: 0},
		},
	}

	return config, expected
}
