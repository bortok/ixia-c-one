package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	location string
	username string
	config   *ssh.ClientConfig
	client   *ssh.Client
}

type DutInterface struct {
	Name     string
	MacAddr  string
	Ipv4Addr string
	Ipv6Addr string
}

type DutIsisNeighbor struct {
	SystemId string
	State    string
}

type IsisStatExpected struct {
	SystemIds   []string
	RoutesCount int
}

func NewSshClient(location string, username string, password string) (*SshClient, error) {
	log.Printf("Creating SSH client for server %s ...\n", location)

	authMethod := []ssh.AuthMethod{}

	if password != "" {
		authMethod = append(authMethod, ssh.Password(password))
	} else {
		authMethod = append(authMethod, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
			return nil, nil
		}))
	}

	sshConfig := ssh.ClientConfig{
		User:            username,
		Auth:            authMethod,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	log.Println("Dialing SSH ...")
	client, err := ssh.Dial("tcp", location, &sshConfig)
	if err != nil {
		return nil, fmt.Errorf("could not dial SSH server %s: %v", location, err)
	}
	return &SshClient{
		location: location,
		username: username,
		config:   &sshConfig,
		client:   client,
	}, nil
}

func (c *SshClient) Close() {
	log.Printf("Closing SSH connection with %s\n", c.location)
	c.client.Close()
}

func (c *SshClient) Exec(cmd string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("could not create ssh session: %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if optsDebug {
		log.Printf("SSH INPUT: %s\n", cmd)
	}
	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("could not execute '%s': %v", cmd, err)
	}

	out := b.String()
	if optsDebug {
		log.Printf("SSH OUTPUT: %s\n", out)
	}
	return out, nil
}

func (c *SshClient) ExecMultiple(cmds []string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("could not create ssh session: %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	for _, cmd := range cmds {
		if optsDebug {
			log.Printf("SSH INPUT: %s\n", cmd)
		}
		if err := session.Run(cmd); err != nil {
			return "", fmt.Errorf("could not execute '%s': %v", cmd, err)
		}
	}

	out := b.String()
	if optsDebug {
		log.Printf("SSH OUTPUT: %s\n", out)
	}
	return out, nil
}

func (c *SshClient) PushDutConfigFile(location string) (string, error) {
	log.Printf("Reading DUT config %s ...", location)
	bytes, err := ioutil.ReadFile(location)
	if err != nil {
		return "", fmt.Errorf("could not read DUT config %s: %v", location, err)
	}

	return c.Exec("enable\nconfig terminal\n" + string(bytes))
}

func (c *SshClient) GetInterface(name string) (*DutInterface, error) {
	ifc := DutInterface{}
	out, err := c.Exec("show interface " + name)
	if err != nil {
		return nil, err
	}

	ifc.Name = name
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "Hardware is Ethernet") {
			for _, word := range strings.Split(line, " ") {
				if strings.Contains(word, ".") {
					ifc.MacAddr = fmt.Sprintf(
						"%s:%s:%s:%s:%s:%s",
						word[0:2], word[2:4], word[5:7], word[7:9], word[10:12],
						word[12:14],
					)
				}
			}
		} else if strings.Contains(line, "Internet address") {
			for _, word := range strings.Split(line, " ") {
				if strings.Contains(word, ".") {
					ifc.Ipv4Addr = strings.Split(word, "/")[0]
				}
			}
		}
	}

	return &ifc, nil
}

func (c *SshClient) CheckRouteInstalled(ipMask string, ifcName string) (bool, error) {
	out, err := c.Exec("show ip route ")
	if err != nil {
		return false, err
	}

	log.Printf("show ip route in DUT:\n %s", out)
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, ifcName) {
			for _, word := range strings.Split(line, " ") {
				if strings.Contains(word, ipMask) {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
