package common

import (
	"strings"
	"testing"
)

// TestConstants_Protocol checks parsing of protocol names
// expectations:
//   - Invalid name raises errors.
//   - tcp and tcps get properly mapped
func TestConstants_Protocol(t *testing.T) {
	t.Parallel()
	var p Protocol
	var err error

	_, err = NormalizeProtocol("")
	if err == nil {
		t.Errorf("should have receive an error for empty string")
	}
	_, err = NormalizeProtocol("XXXx")
	if err == nil {
		t.Errorf("should have receive an error for non-protocol string")
	}

	p, _ = NormalizeProtocol("tcp")
	if p != ProtocolTCP {
		t.Errorf("wrong protocol value for [%s], wanted [%s], got [%s]", "tcp", ProtocolTCP, p)
	}

	p, _ = NormalizeProtocol("tcps")
	if p != ProtocolTCPS {
		t.Errorf("wrong protocol value for [%s], wanted [%s], got [%s]", "tcps", ProtocolTCPS, p)
	}

	p, _ = NormalizeProtocol("TcP")
	if p != ProtocolTCP {
		t.Errorf("wrong protocol value for [%s], wanted [%s], got [%s]", "TcP", ProtocolTCP, p)
	}
	p, _ = NormalizeProtocol("TCPS")
	if p != ProtocolTCPS {
		t.Errorf("wrong protocol value for [%s], wanted [%s], got [%s]", "TCPS", ProtocolTCPS, p)
	}

}

// TestConstants_Protocol checks protocol names  constant mappings
// expectations:
//   - ProtocolTCP and ProtocolTCPS and mapped to tcp and tcps
func TestConstants_ProtocolString(t *testing.T) {
	t.Parallel()
	if strings.Compare(ProtocolTCP.String(), protocolName[ProtocolTCP]) != 0 {
		t.Errorf("wrong protocol string value  for [%s], wanted [%s], got [%s]",
			ProtocolTCP, protocolName[ProtocolTCP],
			ProtocolTCP.String())
	}
	if strings.Compare(ProtocolTCPS.String(), protocolName[ProtocolTCPS]) != 0 {
		t.Errorf("wrong protocol string value  for [%s], wanted [%s], got [%s]",
			ProtocolTCPS, protocolName[ProtocolTCPS],
			ProtocolTCPS.String())
	}
}
