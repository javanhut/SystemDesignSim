package network

import (
	"fmt"
	"net"
	"strings"
)

type VPC struct {
	ID                 string
	Name               string
	Region             string
	CIDR               string
	Subnets            []*Subnet
	InternetGateway    *InternetGateway
	NATGateways        []*NATGateway
	RouteTables        []*RouteTable
	EnableDNS          bool
	EnableDNSHostnames bool
	Tags               map[string]string
	CostPerHour        float64
}

type SubnetType string

const (
	SubnetTypePublic   SubnetType = "public"
	SubnetTypePrivate  SubnetType = "private"
	SubnetTypeDatabase SubnetType = "database"
)

type Subnet struct {
	ID                 string
	Name               string
	VPC                *VPC
	CIDR               string
	AvailabilityZone   string
	Type               SubnetType
	RouteTable         *RouteTable
	AutoAssignPublicIP bool
	AvailableIPs       int
	Tags               map[string]string
}

type InternetGateway struct {
	ID          string
	Name        string
	VPC         *VPC
	Attached    bool
	CostPerHour float64
}

type NATGateway struct {
	ID              string
	Name            string
	Subnet          *Subnet
	ElasticIP       string
	CostPerHour     float64
	DataProcessedGB float64
}

type RouteTable struct {
	ID      string
	Name    string
	VPC     *VPC
	Routes  []*Route
	Subnets []*Subnet
	Main    bool
}

type Route struct {
	DestinationCIDR string
	Target          RouteTarget
	TargetID        string
}

type RouteTarget string

const (
	RouteTargetLocal           RouteTarget = "local"
	RouteTargetInternetGateway RouteTarget = "igw"
	RouteTargetNATGateway      RouteTarget = "nat"
	RouteTargetVPCPeering      RouteTarget = "pcx"
	RouteTargetTransitGateway  RouteTarget = "tgw"
)

func NewVPC(id, name, region, cidr string) (*VPC, error) {
	if !isValidCIDR(cidr) {
		return nil, fmt.Errorf("invalid CIDR: %s", cidr)
	}

	vpc := &VPC{
		ID:                 id,
		Name:               name,
		Region:             region,
		CIDR:               cidr,
		Subnets:            []*Subnet{},
		RouteTables:        []*RouteTable{},
		EnableDNS:          true,
		EnableDNSHostnames: true,
		Tags:               make(map[string]string),
		CostPerHour:        0.0,
	}

	mainRouteTable := &RouteTable{
		ID:      fmt.Sprintf("rtb-%s-main", id),
		Name:    "Main Route Table",
		VPC:     vpc,
		Routes:  []*Route{{DestinationCIDR: cidr, Target: RouteTargetLocal}},
		Subnets: []*Subnet{},
		Main:    true,
	}
	vpc.RouteTables = append(vpc.RouteTables, mainRouteTable)

	return vpc, nil
}

func (vpc *VPC) CreateSubnet(id, name, cidr, az string, subnetType SubnetType) (*Subnet, error) {
	if !isValidCIDR(cidr) {
		return nil, fmt.Errorf("invalid subnet CIDR: %s", cidr)
	}

	if !isSubnetOfVPC(cidr, vpc.CIDR) {
		return nil, fmt.Errorf("subnet CIDR %s is not within VPC CIDR %s", cidr, vpc.CIDR)
	}

	availableIPs := calculateAvailableIPs(cidr)

	subnet := &Subnet{
		ID:                 id,
		Name:               name,
		VPC:                vpc,
		CIDR:               cidr,
		AvailabilityZone:   az,
		Type:               subnetType,
		AutoAssignPublicIP: subnetType == SubnetTypePublic,
		AvailableIPs:       availableIPs,
		Tags:               make(map[string]string),
	}

	mainRouteTable := vpc.GetMainRouteTable()
	if mainRouteTable != nil {
		subnet.RouteTable = mainRouteTable
		mainRouteTable.Subnets = append(mainRouteTable.Subnets, subnet)
	}

	vpc.Subnets = append(vpc.Subnets, subnet)
	return subnet, nil
}

func (vpc *VPC) AttachInternetGateway(id, name string) *InternetGateway {
	igw := &InternetGateway{
		ID:          id,
		Name:        name,
		VPC:         vpc,
		Attached:    true,
		CostPerHour: 0.0,
	}
	vpc.InternetGateway = igw
	return igw
}

func (vpc *VPC) CreateNATGateway(id, name string, subnet *Subnet) (*NATGateway, error) {
	if subnet.Type != SubnetTypePublic {
		return nil, fmt.Errorf("NAT Gateway must be created in a public subnet")
	}

	natGW := &NATGateway{
		ID:              id,
		Name:            name,
		Subnet:          subnet,
		ElasticIP:       fmt.Sprintf("eip-%s", id),
		CostPerHour:     0.045,
		DataProcessedGB: 0.045,
	}

	vpc.NATGateways = append(vpc.NATGateways, natGW)
	return natGW, nil
}

func (vpc *VPC) CreateRouteTable(id, name string) *RouteTable {
	rt := &RouteTable{
		ID:      id,
		Name:    name,
		VPC:     vpc,
		Routes:  []*Route{{DestinationCIDR: vpc.CIDR, Target: RouteTargetLocal}},
		Subnets: []*Subnet{},
		Main:    false,
	}
	vpc.RouteTables = append(vpc.RouteTables, rt)
	return rt
}

func (vpc *VPC) GetMainRouteTable() *RouteTable {
	for _, rt := range vpc.RouteTables {
		if rt.Main {
			return rt
		}
	}
	return nil
}

func (vpc *VPC) GetSubnet(id string) *Subnet {
	for _, subnet := range vpc.Subnets {
		if subnet.ID == id {
			return subnet
		}
	}
	return nil
}

func (vpc *VPC) GetPublicSubnets() []*Subnet {
	subnets := []*Subnet{}
	for _, subnet := range vpc.Subnets {
		if subnet.Type == SubnetTypePublic {
			subnets = append(subnets, subnet)
		}
	}
	return subnets
}

func (vpc *VPC) GetPrivateSubnets() []*Subnet {
	subnets := []*Subnet{}
	for _, subnet := range vpc.Subnets {
		if subnet.Type == SubnetTypePrivate {
			subnets = append(subnets, subnet)
		}
	}
	return subnets
}

func (vpc *VPC) CalculateTotalCost() float64 {
	totalCost := vpc.CostPerHour

	for _, natGW := range vpc.NATGateways {
		totalCost += natGW.CostPerHour
	}

	return totalCost
}

func (rt *RouteTable) AddRoute(destinationCIDR string, target RouteTarget, targetID string) error {
	for _, route := range rt.Routes {
		if route.DestinationCIDR == destinationCIDR {
			return fmt.Errorf("route for %s already exists", destinationCIDR)
		}
	}

	route := &Route{
		DestinationCIDR: destinationCIDR,
		Target:          target,
		TargetID:        targetID,
	}
	rt.Routes = append(rt.Routes, route)
	return nil
}

func (rt *RouteTable) AssociateSubnet(subnet *Subnet) {
	if subnet.RouteTable != nil {
		for i, s := range subnet.RouteTable.Subnets {
			if s.ID == subnet.ID {
				subnet.RouteTable.Subnets = append(
					subnet.RouteTable.Subnets[:i],
					subnet.RouteTable.Subnets[i+1:]...,
				)
				break
			}
		}
	}

	subnet.RouteTable = rt
	rt.Subnets = append(rt.Subnets, subnet)
}

func (subnet *Subnet) MakePublic(igw *InternetGateway) error {
	if igw == nil || !igw.Attached {
		return fmt.Errorf("Internet Gateway not attached to VPC")
	}

	subnet.Type = SubnetTypePublic
	subnet.AutoAssignPublicIP = true

	if subnet.RouteTable == nil {
		rt := subnet.VPC.CreateRouteTable(
			fmt.Sprintf("rtb-%s-public", subnet.ID),
			fmt.Sprintf("%s Public Route Table", subnet.Name),
		)
		rt.AssociateSubnet(subnet)
	}

	subnet.RouteTable.AddRoute("0.0.0.0/0", RouteTargetInternetGateway, igw.ID)

	return nil
}

func (subnet *Subnet) MakePrivate(natGW *NATGateway) error {
	subnet.Type = SubnetTypePrivate
	subnet.AutoAssignPublicIP = false

	if subnet.RouteTable == nil {
		rt := subnet.VPC.CreateRouteTable(
			fmt.Sprintf("rtb-%s-private", subnet.ID),
			fmt.Sprintf("%s Private Route Table", subnet.Name),
		)
		rt.AssociateSubnet(subnet)
	}

	if natGW != nil {
		subnet.RouteTable.AddRoute("0.0.0.0/0", RouteTargetNATGateway, natGW.ID)
	}

	return nil
}

func isValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

func isSubnetOfVPC(subnetCIDR, vpcCIDR string) bool {
	_, vpcNet, err := net.ParseCIDR(vpcCIDR)
	if err != nil {
		return false
	}

	subnetIP, _, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return false
	}

	return vpcNet.Contains(subnetIP)
}

func calculateAvailableIPs(cidr string) int {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0
	}

	ones, bits := ipNet.Mask.Size()
	hostBits := bits - ones
	totalIPs := 1 << uint(hostBits)

	return totalIPs - 5
}

type VPCPeering struct {
	ID           string
	Name         string
	RequesterVPC *VPC
	AccepterVPC  *VPC
	Status       PeeringStatus
	CostPerGB    float64
}

type PeeringStatus string

const (
	PeeringStatusPending  PeeringStatus = "pending"
	PeeringStatusActive   PeeringStatus = "active"
	PeeringStatusRejected PeeringStatus = "rejected"
	PeeringStatusDeleted  PeeringStatus = "deleted"
)

func NewVPCPeering(id, name string, requester, accepter *VPC) *VPCPeering {
	return &VPCPeering{
		ID:           id,
		Name:         name,
		RequesterVPC: requester,
		AccepterVPC:  accepter,
		Status:       PeeringStatusPending,
		CostPerGB:    0.01,
	}
}

func (p *VPCPeering) Accept() {
	p.Status = PeeringStatusActive
}

func (p *VPCPeering) Reject() {
	p.Status = PeeringStatusRejected
}

type VPCPreset struct {
	Name        string
	Description string
	CIDR        string
	Subnets     []SubnetPreset
}

type SubnetPreset struct {
	Name string
	CIDR string
	Type SubnetType
}

var VPCPresets = map[string]*VPCPreset{
	"single-az": {
		Name:        "Single AZ",
		Description: "Basic VPC with public and private subnets in one AZ",
		CIDR:        "10.0.0.0/16",
		Subnets: []SubnetPreset{
			{Name: "Public Subnet", CIDR: "10.0.1.0/24", Type: SubnetTypePublic},
			{Name: "Private Subnet", CIDR: "10.0.2.0/24", Type: SubnetTypePrivate},
		},
	},
	"multi-az": {
		Name:        "Multi AZ",
		Description: "High availability VPC with subnets across 3 AZs",
		CIDR:        "10.0.0.0/16",
		Subnets: []SubnetPreset{
			{Name: "Public Subnet AZ1", CIDR: "10.0.1.0/24", Type: SubnetTypePublic},
			{Name: "Public Subnet AZ2", CIDR: "10.0.2.0/24", Type: SubnetTypePublic},
			{Name: "Public Subnet AZ3", CIDR: "10.0.3.0/24", Type: SubnetTypePublic},
			{Name: "Private Subnet AZ1", CIDR: "10.0.11.0/24", Type: SubnetTypePrivate},
			{Name: "Private Subnet AZ2", CIDR: "10.0.12.0/24", Type: SubnetTypePrivate},
			{Name: "Private Subnet AZ3", CIDR: "10.0.13.0/24", Type: SubnetTypePrivate},
		},
	},
	"three-tier": {
		Name:        "Three Tier",
		Description: "Web, app, and database tiers across 2 AZs",
		CIDR:        "10.0.0.0/16",
		Subnets: []SubnetPreset{
			{Name: "Web Subnet AZ1", CIDR: "10.0.1.0/24", Type: SubnetTypePublic},
			{Name: "Web Subnet AZ2", CIDR: "10.0.2.0/24", Type: SubnetTypePublic},
			{Name: "App Subnet AZ1", CIDR: "10.0.11.0/24", Type: SubnetTypePrivate},
			{Name: "App Subnet AZ2", CIDR: "10.0.12.0/24", Type: SubnetTypePrivate},
			{Name: "DB Subnet AZ1", CIDR: "10.0.21.0/24", Type: SubnetTypeDatabase},
			{Name: "DB Subnet AZ2", CIDR: "10.0.22.0/24", Type: SubnetTypeDatabase},
		},
	},
}

func GetVPCPreset(name string) *VPCPreset {
	if preset, exists := VPCPresets[name]; exists {
		return preset
	}
	return VPCPresets["single-az"]
}

func GetVPCPresetNames() []string {
	return []string{"single-az", "multi-az", "three-tier"}
}

func CreateVPCFromPreset(vpcID, vpcName, region string, presetName string, azs []string) (*VPC, error) {
	preset := GetVPCPreset(presetName)

	vpc, err := NewVPC(vpcID, vpcName, region, preset.CIDR)
	if err != nil {
		return nil, err
	}

	igw := vpc.AttachInternetGateway(fmt.Sprintf("igw-%s", vpcID), fmt.Sprintf("%s IGW", vpcName))

	azIndex := 0
	for i, subnetPreset := range preset.Subnets {
		if len(azs) == 0 {
			azs = []string{"a"}
		}

		az := azs[azIndex%len(azs)]
		subnetID := fmt.Sprintf("subnet-%s-%d", vpcID, i)

		subnet, err := vpc.CreateSubnet(subnetID, subnetPreset.Name, subnetPreset.CIDR, az, subnetPreset.Type)
		if err != nil {
			return nil, err
		}

		if subnetPreset.Type == SubnetTypePublic {
			subnet.MakePublic(igw)
		}

		if strings.Contains(strings.ToLower(subnetPreset.Name), "az2") ||
			strings.Contains(strings.ToLower(subnetPreset.Name), "az3") {
			azIndex++
		}
	}

	if len(vpc.GetPublicSubnets()) > 0 {
		publicSubnet := vpc.GetPublicSubnets()[0]
		natGW, _ := vpc.CreateNATGateway(fmt.Sprintf("nat-%s", vpcID), "NAT Gateway", publicSubnet)

		for _, subnet := range vpc.GetPrivateSubnets() {
			subnet.MakePrivate(natGW)
		}
	}

	return vpc, nil
}
