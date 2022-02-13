/* Test L2 Traffic
Topology:
IXIA (40.40.40.0/24) -----> ARISTA ------> IXIA (50.50.50.0/24)
Flows:
- v4: 40.40.40.1 -> 50.50.50.1
*/

package tests

import (
	"testing"

	"github.com/open-traffic-generator/ixia-c-one/tests/helpers"
	"github.com/open-traffic-generator/snappi/gosnappi"
)

// This test is verified with ceos which by default acts as a switch and hence no configuration
// on DUT for L2 traffic test
func TestL2Traffic(t *testing.T) {
	client, err := helpers.NewClient(otgHttpLocation)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	config, expected := trafficConfigL2(client)

	if err := client.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	if err := client.StartTransmit(nil); err != nil {
		t.Fatal(err)
	}

	helpers.WaitFor(t, func() (bool, error) { return client.FlowMetricsOk(expected) }, nil)

}

func trafficConfigL2(client *helpers.ApiClient) (gosnappi.Config, helpers.ExpectedState) {
	config := client.Api().NewConfig()

	port1 := config.Ports().Add().SetName("ixia-c-port1").SetLocation(otgPort1Location)
	port2 := config.Ports().Add().SetName("ixia-c-port2").SetLocation(otgPort2Location)

	// OTG traffic configuration
	f1 := config.Flows().Add().SetName("p1.v4.p2")
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

	expected := helpers.ExpectedState{
		Flow: map[string]helpers.ExpectedFlowMetrics{
			f1.Name(): {FramesRx: 1000, FramesRxRate: 0},
		},
	}

	return config, expected
}
