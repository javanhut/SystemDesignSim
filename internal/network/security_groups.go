package network

import (
	"fmt"
	"strconv"
	"strings"
)

type SecurityGroup struct {
	ID           string
	Name         string
	Description  string
	VPC          *VPC
	IngressRules []*SecurityRule
	EgressRules  []*SecurityRule
	Tags         map[string]string
}

type SecurityRule struct {
	ID          string
	Protocol    Protocol
	PortRange   PortRange
	Source      string
	Description string
	RuleType    RuleType
}

type Protocol string

const (
	ProtocolTCP  Protocol = "tcp"
	ProtocolUDP  Protocol = "udp"
	ProtocolICMP Protocol = "icmp"
	ProtocolAll  Protocol = "-1"
)

type RuleType string

const (
	RuleTypeIngress RuleType = "ingress"
	RuleTypeEgress  RuleType = "egress"
)

type PortRange struct {
	From int
	To   int
}

func NewSecurityGroup(id, name, description string, vpc *VPC) *SecurityGroup {
	sg := &SecurityGroup{
		ID:           id,
		Name:         name,
		Description:  description,
		VPC:          vpc,
		IngressRules: []*SecurityRule{},
		EgressRules: []*SecurityRule{
			{
				ID:          fmt.Sprintf("%s-egress-all", id),
				Protocol:    ProtocolAll,
				PortRange:   PortRange{From: 0, To: 65535},
				Source:      "0.0.0.0/0",
				Description: "Allow all outbound traffic",
				RuleType:    RuleTypeEgress,
			},
		},
		Tags: make(map[string]string),
	}
	return sg
}

func (sg *SecurityGroup) AddIngressRule(protocol Protocol, fromPort, toPort int, source, description string) *SecurityRule {
	rule := &SecurityRule{
		ID:          fmt.Sprintf("%s-ingress-%d", sg.ID, len(sg.IngressRules)),
		Protocol:    protocol,
		PortRange:   PortRange{From: fromPort, To: toPort},
		Source:      source,
		Description: description,
		RuleType:    RuleTypeIngress,
	}
	sg.IngressRules = append(sg.IngressRules, rule)
	return rule
}

func (sg *SecurityGroup) AddEgressRule(protocol Protocol, fromPort, toPort int, destination, description string) *SecurityRule {
	rule := &SecurityRule{
		ID:          fmt.Sprintf("%s-egress-%d", sg.ID, len(sg.EgressRules)),
		Protocol:    protocol,
		PortRange:   PortRange{From: fromPort, To: toPort},
		Source:      destination,
		Description: description,
		RuleType:    RuleTypeEgress,
	}
	sg.EgressRules = append(sg.EgressRules, rule)
	return rule
}

func (sg *SecurityGroup) RemoveIngressRule(ruleID string) error {
	for i, rule := range sg.IngressRules {
		if rule.ID == ruleID {
			sg.IngressRules = append(sg.IngressRules[:i], sg.IngressRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("ingress rule %s not found", ruleID)
}

func (sg *SecurityGroup) RemoveEgressRule(ruleID string) error {
	for i, rule := range sg.EgressRules {
		if rule.ID == ruleID {
			sg.EgressRules = append(sg.EgressRules[:i], sg.EgressRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("egress rule %s not found", ruleID)
}

func (sg *SecurityGroup) AllowHTTP(source string) *SecurityRule {
	return sg.AddIngressRule(ProtocolTCP, 80, 80, source, "Allow HTTP traffic")
}

func (sg *SecurityGroup) AllowHTTPS(source string) *SecurityRule {
	return sg.AddIngressRule(ProtocolTCP, 443, 443, source, "Allow HTTPS traffic")
}

func (sg *SecurityGroup) AllowSSH(source string) *SecurityRule {
	return sg.AddIngressRule(ProtocolTCP, 22, 22, source, "Allow SSH access")
}

func (sg *SecurityGroup) AllowMySQL(source string) *SecurityRule {
	return sg.AddIngressRule(ProtocolTCP, 3306, 3306, source, "Allow MySQL traffic")
}

func (sg *SecurityGroup) AllowPostgreSQL(source string) *SecurityRule {
	return sg.AddIngressRule(ProtocolTCP, 5432, 5432, source, "Allow PostgreSQL traffic")
}

func (sg *SecurityGroup) AllowRedis(source string) *SecurityRule {
	return sg.AddIngressRule(ProtocolTCP, 6379, 6379, source, "Allow Redis traffic")
}

func (sg *SecurityGroup) AllowMemcached(source string) *SecurityRule {
	return sg.AddIngressRule(ProtocolTCP, 11211, 11211, source, "Allow Memcached traffic")
}

func (sg *SecurityGroup) AllowCustomPort(port int, protocol Protocol, source, description string) *SecurityRule {
	return sg.AddIngressRule(protocol, port, port, source, description)
}

func (sg *SecurityGroup) AllowPortRange(fromPort, toPort int, protocol Protocol, source, description string) *SecurityRule {
	return sg.AddIngressRule(protocol, fromPort, toPort, source, description)
}

func (sg *SecurityGroup) IsTrafficAllowed(protocol Protocol, port int, sourceIP string) bool {
	for _, rule := range sg.IngressRules {
		if rule.Protocol == ProtocolAll || rule.Protocol == protocol {
			if port >= rule.PortRange.From && port <= rule.PortRange.To {
				if matchesSource(sourceIP, rule.Source) {
					return true
				}
			}
		}
	}
	return false
}

func matchesSource(ip, source string) bool {
	if source == "0.0.0.0/0" || source == "::/0" {
		return true
	}

	if strings.HasPrefix(source, "sg-") {
		return false
	}

	if strings.Contains(source, "/") {
		return true
	}

	return ip == source
}

type NetworkACL struct {
	ID           string
	Name         string
	VPC          *VPC
	Subnets      []*Subnet
	IngressRules []*ACLRule
	EgressRules  []*ACLRule
}

type ACLRule struct {
	RuleNumber  int
	Protocol    Protocol
	PortRange   PortRange
	Source      string
	Action      ACLAction
	Description string
}

type ACLAction string

const (
	ACLActionAllow ACLAction = "allow"
	ACLActionDeny  ACLAction = "deny"
)

func NewNetworkACL(id, name string, vpc *VPC) *NetworkACL {
	return &NetworkACL{
		ID:           id,
		Name:         name,
		VPC:          vpc,
		Subnets:      []*Subnet{},
		IngressRules: []*ACLRule{},
		EgressRules:  []*ACLRule{},
	}
}

func (nacl *NetworkACL) AddIngressRule(ruleNumber int, protocol Protocol, fromPort, toPort int, source string, action ACLAction, description string) *ACLRule {
	rule := &ACLRule{
		RuleNumber:  ruleNumber,
		Protocol:    protocol,
		PortRange:   PortRange{From: fromPort, To: toPort},
		Source:      source,
		Action:      action,
		Description: description,
	}
	nacl.IngressRules = append(nacl.IngressRules, rule)
	return rule
}

func (nacl *NetworkACL) AddEgressRule(ruleNumber int, protocol Protocol, fromPort, toPort int, destination string, action ACLAction, description string) *ACLRule {
	rule := &ACLRule{
		RuleNumber:  ruleNumber,
		Protocol:    protocol,
		PortRange:   PortRange{From: fromPort, To: toPort},
		Source:      destination,
		Action:      action,
		Description: description,
	}
	nacl.EgressRules = append(nacl.EgressRules, rule)
	return rule
}

func (nacl *NetworkACL) AssociateSubnet(subnet *Subnet) {
	nacl.Subnets = append(nacl.Subnets, subnet)
}

type SecurityGroupPreset struct {
	Name        string
	Description string
	Rules       []PresetRule
}

type PresetRule struct {
	Protocol    Protocol
	FromPort    int
	ToPort      int
	Source      string
	Description string
}

var SecurityGroupPresets = map[string]*SecurityGroupPreset{
	"web-server": {
		Name:        "Web Server",
		Description: "Allow HTTP, HTTPS from anywhere",
		Rules: []PresetRule{
			{Protocol: ProtocolTCP, FromPort: 80, ToPort: 80, Source: "0.0.0.0/0", Description: "HTTP from anywhere"},
			{Protocol: ProtocolTCP, FromPort: 443, ToPort: 443, Source: "0.0.0.0/0", Description: "HTTPS from anywhere"},
		},
	},
	"app-server": {
		Name:        "Application Server",
		Description: "Allow custom app ports from load balancer",
		Rules: []PresetRule{
			{Protocol: ProtocolTCP, FromPort: 8080, ToPort: 8080, Source: "10.0.0.0/16", Description: "App port from VPC"},
			{Protocol: ProtocolTCP, FromPort: 3000, ToPort: 3000, Source: "10.0.0.0/16", Description: "Node.js from VPC"},
		},
	},
	"database": {
		Name:        "Database",
		Description: "Allow database ports from app tier",
		Rules: []PresetRule{
			{Protocol: ProtocolTCP, FromPort: 3306, ToPort: 3306, Source: "10.0.0.0/16", Description: "MySQL from VPC"},
			{Protocol: ProtocolTCP, FromPort: 5432, ToPort: 5432, Source: "10.0.0.0/16", Description: "PostgreSQL from VPC"},
		},
	},
	"cache": {
		Name:        "Cache",
		Description: "Allow cache ports from app tier",
		Rules: []PresetRule{
			{Protocol: ProtocolTCP, FromPort: 6379, ToPort: 6379, Source: "10.0.0.0/16", Description: "Redis from VPC"},
			{Protocol: ProtocolTCP, FromPort: 11211, ToPort: 11211, Source: "10.0.0.0/16", Description: "Memcached from VPC"},
		},
	},
	"load-balancer": {
		Name:        "Load Balancer",
		Description: "Allow HTTP/HTTPS from anywhere",
		Rules: []PresetRule{
			{Protocol: ProtocolTCP, FromPort: 80, ToPort: 80, Source: "0.0.0.0/0", Description: "HTTP from anywhere"},
			{Protocol: ProtocolTCP, FromPort: 443, ToPort: 443, Source: "0.0.0.0/0", Description: "HTTPS from anywhere"},
		},
	},
}

func GetSecurityGroupPreset(name string) *SecurityGroupPreset {
	if preset, exists := SecurityGroupPresets[name]; exists {
		return preset
	}
	return SecurityGroupPresets["web-server"]
}

func GetSecurityGroupPresetNames() []string {
	return []string{"web-server", "app-server", "database", "cache", "load-balancer"}
}

func CreateSecurityGroupFromPreset(id, name string, vpc *VPC, presetName string) *SecurityGroup {
	preset := GetSecurityGroupPreset(presetName)

	sg := NewSecurityGroup(id, name, preset.Description, vpc)

	for _, rule := range preset.Rules {
		sg.AddIngressRule(rule.Protocol, rule.FromPort, rule.ToPort, rule.Source, rule.Description)
	}

	return sg
}

func ParsePort(portStr string) (int, error) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("invalid port: %s", portStr)
	}
	if port < 0 || port > 65535 {
		return 0, fmt.Errorf("port out of range: %d", port)
	}
	return port, nil
}

func ParsePortRange(rangeStr string) (PortRange, error) {
	parts := strings.Split(rangeStr, "-")
	if len(parts) == 1 {
		port, err := ParsePort(parts[0])
		if err != nil {
			return PortRange{}, err
		}
		return PortRange{From: port, To: port}, nil
	}

	if len(parts) == 2 {
		fromPort, err := ParsePort(parts[0])
		if err != nil {
			return PortRange{}, err
		}
		toPort, err := ParsePort(parts[1])
		if err != nil {
			return PortRange{}, err
		}
		if fromPort > toPort {
			return PortRange{}, fmt.Errorf("invalid port range: from port greater than to port")
		}
		return PortRange{From: fromPort, To: toPort}, nil
	}

	return PortRange{}, fmt.Errorf("invalid port range format: %s", rangeStr)
}

func (pr PortRange) String() string {
	if pr.From == pr.To {
		return fmt.Sprintf("%d", pr.From)
	}
	return fmt.Sprintf("%d-%d", pr.From, pr.To)
}

func (pr PortRange) Contains(port int) bool {
	return port >= pr.From && port <= pr.To
}
