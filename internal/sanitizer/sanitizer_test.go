package sanitizer

import (
	"strings"
	"testing"
)

func TestSanitizePassword(t *testing.T) {
	cfg := SanitizeConfig{
		Enabled:           true,
		MaskPasswords:     true,
		MaskSNMPCommunity: false,
		MaskEnableSecret:  false,
	}
	s := NewSanitizer(cfg)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "cisco password direct",
			input:    "password mysecret123",
			expected: "password ********",
		},
		{
			name:     "passwd field",
			input:    "passwd admin123",
			expected: "passwd ********",
		},
		{
			name:     "case insensitive",
			input:    "Password MySecret",
			expected: "Password ********",
		},
		{
			name:     "h3c password with modifier",
			input:    "password simple h3cpass",
			expected: "password ******** ********",
		},
		{
			name:     "huawei password with modifier",
			input:    "password cipher huaweipass",
			expected: "password ******** ********",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.Sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("Sanitize() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizeSNMPCommunity(t *testing.T) {
	cfg := SanitizeConfig{
		Enabled:           true,
		MaskPasswords:     false,
		MaskSNMPCommunity: true,
		MaskEnableSecret:  false,
	}
	s := NewSanitizer(cfg)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "cisco snmp community ro",
			input:    "snmp-server community public RO",
			expected: "snmp-server community ******** RO",
		},
		{
			name:     "cisco snmp community rw",
			input:    "snmp-server community private RW",
			expected: "snmp-server community ******** RW",
		},
		{
			name:     "case insensitive",
			input:    "SNMP-server Community teststr RO",
			expected: "SNMP-server Community ******** RO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.Sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("Sanitize() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizeEnableSecret(t *testing.T) {
	cfg := SanitizeConfig{
		Enabled:           true,
		MaskPasswords:     false,
		MaskSNMPCommunity: false,
		MaskEnableSecret:  true,
	}
	s := NewSanitizer(cfg)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "enable secret with type",
			input:    "enable secret 5 $1$xyz",
			expected: "enable secret ******** ********",
		},
		{
			name:     "enable password direct",
			input:    "enable password cisco",
			expected: "enable password ********",
		},
		{
			name:     "case insensitive",
			input:    "Enable Secret mykey",
			expected: "Enable Secret ********",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.Sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("Sanitize() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizeMultiVendor(t *testing.T) {
	cfg := SanitizeConfig{
		Enabled:           true,
		MaskPasswords:     true,
		MaskSNMPCommunity: true,
		MaskEnableSecret:  true,
	}
	s := NewSanitizer(cfg)

	ciscoConfig := `!
hostname Router1
enable secret 5 $1$xyz
password cisco123
snmp-server community public RO
interface GigabitEthernet0/0
 ip address 192.168.1.1 255.255.255.0
!`

	result := s.Sanitize(ciscoConfig)
	if strings.Contains(result, "cisco123") {
		t.Error("cisco password not masked")
	}
	if strings.Contains(result, "$1$xyz") {
		t.Error("enable secret not masked")
	}
	if strings.Contains(result, "public") && !strings.Contains(result, "snmp-server community ********") {
		t.Error("snmp community not masked")
	}

	h3cConfig := `#
 sysname Switch1
 password simple h3cpass
 snmp-server community h3cro RO
#`

	result = s.Sanitize(h3cConfig)
	if strings.Contains(result, "h3cpass") {
		t.Error("h3c password not masked")
	}

	huaweiConfig := `#
 sysname Huawei1
 password cipher huaweipass
 snmp-server community huawei RO
#`

	result = s.Sanitize(huaweiConfig)
	if strings.Contains(result, "huaweipass") {
		t.Error("huawei password not masked")
	}

	ruijieConfig := `!
hostname Ruijie1
password ruijiepass
snmp-server community ruijie RO
!`

	result = s.Sanitize(ruijieConfig)
	if strings.Contains(result, "ruijiepass") {
		t.Error("ruijie password not masked")
	}
}

func TestSanitizeDisabled(t *testing.T) {
	cfg := SanitizeConfig{
		Enabled:           false,
		MaskPasswords:     true,
		MaskSNMPCommunity: true,
		MaskEnableSecret:  true,
	}
	s := NewSanitizer(cfg)
	input := "password secret123"
	result := s.Sanitize(input)
	if result != input {
		t.Errorf("Sanitize with disabled should return original, got %q", result)
	}
}

func TestSanitizeNoPatterns(t *testing.T) {
	cfg := SanitizeConfig{
		Enabled:           true,
		MaskPasswords:     false,
		MaskSNMPCommunity: false,
		MaskEnableSecret:  false,
	}
	s := NewSanitizer(cfg)
	input := "password secret123"
	result := s.Sanitize(input)
	if result != input {
		t.Errorf("Sanitize with no patterns should return original, got %q", result)
	}
}
