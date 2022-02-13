/* Test L3 Traffic
Topology:
IXIA (11.11.11.2/24) -----> (11.11.11.1)ARISTA(12.12.12.1) ------> IXIA (12.12.12.2/24)
Flows:
- tcp: 40.40.40.1 -> 50.50.50.1+
*/

package tests

import (
	"testing"

	"github.com/open-traffic-generator/ixia-c-one/tests/helpers"
	"github.com/open-traffic-generator/snappi/gosnappi"
)

func TestL3Traffic(t *testing.T) {
	client, err := helpers.NewClient(otgHttpLocation)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	config, expected := trafficConfigL3(client)

	// DUT configurationS
	dut, err := helpers.NewSshClient(dutSshLocation, "admin", "")
	if err != nil {
		t.Fatal(err)
	}
	defer dut.Close()

	if _, err := dut.PushDutConfigFile("configs/l3_traffic/set_dut.txt"); err != nil {
		t.Fatal(err)
	}
	defer dut.PushDutConfigFile("configs/l3_traffic/unset_dut.txt")

	// Set the eth destination in configuration from retrieved mac of DUT interface interfacing ixiac Tx port
	ifc, err := dut.GetInterface("Ethernet1")
	if err != nil {
		t.Error(err)
	}
	config.Flows().Items()[0].Packet().Items()[0].Ethernet().Dst().SetValue(ifc.MacAddr)

	// Set ipv4 address in ixiac-te interface of Rx ports
	// (as Traffic-Engine doesn't reply to ARP, by doing this linux will reply to the ARP.)
	if err := helpers.SetTrafficEndPointV4Config(
		"clab-ixia-c-ixia-c-one",
		"eth2",
		"12.12.12.2",
		24); err != nil {
		t.Fatal(err)
	}
	defer helpers.UnSetTrafficEndPointV4Config("clab-ixia-c-ixia-c-one", "eth2", "12.12.12.2", 24)

	// Send OTG traffic
	if err := client.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	if err := client.StartTransmit(nil); err != nil {
		t.Fatal(err)
	}

	helpers.WaitFor(t, func() (bool, error) { return client.FlowMetricsOk(expected) }, nil)

}

func trafficConfigL3(client *helpers.ApiClient) (gosnappi.Config, helpers.ExpectedState) {
	config := client.Api().NewConfig()

	port1 := config.Ports().Add().SetName("ixia-c-port1").SetLocation(otgPort1Location)
	port2 := config.Ports().Add().SetName("ixia-c-port2").SetLocation(otgPort2Location)

	// OTG traffic configuration
	f1 := config.Flows().Add().SetName("p1.tcp.v4.p2")
	f1.Metrics().SetEnable(true)
	f1.TxRx().Port().
		SetTxName(port1.Name()).
		SetRxName(port2.Name())
	f1.Size().SetFixed(512)
	f1.Rate().SetPps(500)
	f1.Duration().FixedPackets().SetPackets(1000)
	e1 := f1.Packet().Add().Ethernet()
	e1.Src().SetValue("00:00:00:00:00:AA")
	e1.Dst().SetValue("00:00:00:00:00:BB")
	v4 := f1.Packet().Add().Ipv4()
	v4.Src().SetValue("40.40.40.1")
	v4.Dst().Increment().SetStart("50.50.50.1").SetStep("0.0.0.1").SetCount(5)
	tc := f1.Packet().Add().Tcp()
	tc.SrcPort().SetValue(3250)
	tc.DstPort().Decrement().SetStart(8070).SetStep(2).SetCount(10)

	expected := helpers.ExpectedState{
		Flow: map[string]helpers.ExpectedFlowMetrics{
			f1.Name(): {FramesRx: 1000, FramesRxRate: 0},
		},
	}

	return config, expected
}
