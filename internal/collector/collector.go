package collector

import (
	"fmt"
	"time"

	"github.com/netconfighub/netconfighub/internal/model"
)

type UnifiedCollector struct {
	sshClient    *SSHCollector
	telnetClient *TelnetCollector
}

func NewCollector() *UnifiedCollector {
	return &UnifiedCollector{
		sshClient:    NewSSHCollector(10 * time.Second),
		telnetClient: NewTelnetCollector(10 * time.Second),
	}
}

func (c *UnifiedCollector) Collect(device *model.Device) (string, error) {
	switch device.Protocol {
	case model.ConnProtocolSSH:
		return c.sshClient.Collect(*device)
	case model.ConnProtocolTelnet:
		return c.telnetClient.Collect(*device)
	default:
		return "", fmt.Errorf("不支持的连接协议: %s", device.Protocol)
	}
}
