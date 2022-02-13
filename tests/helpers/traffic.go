package helpers

import (
	"fmt"
)

func pushConfigToEndPoint(containerName string, configCmd string) error {
	bashcmd := fmt.Sprintf("docker exec %s bash %s", containerName, configCmd)
	//cmd := exec.Command("/bin/sh", "-c", bashcmd)
	fmt.Println("CMD fired: ", bashcmd)
	// _, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return err
	// }
	return nil
}

func SetTrafficEndPointV4Config(containerName string, ifcName string, ipv4Addr string, mask int) error {
	cfgCmd := fmt.Sprintf("set ipv4 %s %s %d", ifcName, ipv4Addr, mask)
	if err := pushConfigToEndPoint(containerName, cfgCmd); err != nil {
		return err
	}

	return nil
}

func UnSetTrafficEndPointV4Config(containerName string, ifcName string, ipv4Addr string, mask int) error {
	cfgCmd := fmt.Sprintf("unset ipv4 %s %s %d", ifcName, ipv4Addr, mask)
	if err := pushConfigToEndPoint(containerName, cfgCmd); err != nil {
		return err
	}

	return nil
}
