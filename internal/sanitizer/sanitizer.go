package sanitizer

import (
	"regexp"
)

type Sanitizer struct {
	patterns []*SanitizePattern
}

type SanitizePattern struct {
	Name    string
	Regex   *regexp.Regexp
	Replace string
}

type SanitizeConfig struct {
	Enabled           bool `yaml:"enabled" json:"enabled"`
	MaskPasswords     bool `yaml:"mask_passwords" json:"mask_passwords"`
	MaskSNMPCommunity bool `yaml:"mask_snmp_community" json:"mask_snmp_community"`
	MaskEnableSecret  bool `yaml:"mask_enable_secret" json:"mask_enable_secret"`
}

func NewSanitizer(config SanitizeConfig) *Sanitizer {
	s := &Sanitizer{}
	if !config.Enabled {
		return s
	}
	if config.MaskEnableSecret {
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "enable_secret_with_type",
			Regex:   regexp.MustCompile(`(?i)(enable[ \t]+(?:secret|password))[ \t]+\S+[ \t]+(\S+)`),
			Replace: "${1} ******** ********",
		})
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "enable_secret_direct",
			Regex:   regexp.MustCompile(`(?i)(enable[ \t]+(?:secret|password))[ \t]+(\S+)`),
			Replace: "${1} ********",
		})
	}
	if config.MaskPasswords {
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "password_with_modifier",
			Regex:   regexp.MustCompile(`(?i)(password|passwd)[ \t]+\S+[ \t]+(\S+)`),
			Replace: "${1} ******** ********",
		})
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "password_direct",
			Regex:   regexp.MustCompile(`(?i)(password|passwd)[ \t]+(\S+)`),
			Replace: "${1} ********",
		})
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "h3c_super_password",
			Regex:   regexp.MustCompile(`(?i)(super[ \t]+password)[ \t]+(\S+)`),
			Replace: "${1} ********",
		})
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "huawei_local_user_password",
			Regex:   regexp.MustCompile(`(?i)(local-user[ \t]+\S+[ \t]+password)[ \t]+(?:cipher|simple)[ \t]+(\S+)`),
			Replace: "${1} ******** ********",
		})
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "h3c_password_simple_cipher",
			Regex:   regexp.MustCompile(`(?i)(password)[ \t]+(?:simple|cipher)[ \t]+(\S+)`),
			Replace: "${1} ******** ********",
		})
	}
	if config.MaskSNMPCommunity {
		s.patterns = append(s.patterns, &SanitizePattern{
			Name:    "snmp_community",
			Regex:   regexp.MustCompile(`(?i)(snmp-server[ \t]+community)[ \t]+(\S+)`),
			Replace: "${1} ********",
		})
	}
	return s
}

func (s *Sanitizer) Sanitize(content string) string {
	result := content
	for _, p := range s.patterns {
		result = p.Regex.ReplaceAllString(result, p.Replace)
	}
	return result
}
