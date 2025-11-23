package network

import (
	"fmt"
	"strings"
	"time"
)

type DNSRecord struct {
	Name        string
	Type        RecordType
	Value       string
	TTL         int
	Region      string
	HealthCheck *HealthCheck
}

type RecordType string

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeTXT   RecordType = "TXT"
	RecordTypeNS    RecordType = "NS"
	RecordTypeSOA   RecordType = "SOA"
	RecordTypeSRV   RecordType = "SRV"
)

type HealthCheck struct {
	ID              string
	Protocol        string
	Port            int
	Path            string
	IntervalSeconds int
	TimeoutSeconds  int
	Healthy         bool
	LastCheck       time.Time
}

type HostedZone struct {
	ID           string
	Name         string
	Records      []*DNSRecord
	Private      bool
	VPC          *VPC
	CostPerZone  float64
	CostPerQuery float64
}

func NewHostedZone(id, name string, private bool, vpc *VPC) *HostedZone {
	return &HostedZone{
		ID:           id,
		Name:         name,
		Records:      []*DNSRecord{},
		Private:      private,
		VPC:          vpc,
		CostPerZone:  0.50,
		CostPerQuery: 0.00000040,
	}
}

func (hz *HostedZone) AddRecord(name string, recordType RecordType, value string, ttl int) *DNSRecord {
	record := &DNSRecord{
		Name:   name,
		Type:   recordType,
		Value:  value,
		TTL:    ttl,
		Region: "",
	}
	hz.Records = append(hz.Records, record)
	return record
}

func (hz *HostedZone) AddARecord(name, ip string, ttl int) *DNSRecord {
	return hz.AddRecord(name, RecordTypeA, ip, ttl)
}

func (hz *HostedZone) AddCNAMERecord(name, target string, ttl int) *DNSRecord {
	return hz.AddRecord(name, RecordTypeCNAME, target, ttl)
}

func (hz *HostedZone) GetRecord(name string, recordType RecordType) *DNSRecord {
	for _, record := range hz.Records {
		if record.Name == name && record.Type == recordType {
			return record
		}
	}
	return nil
}

func (hz *HostedZone) DeleteRecord(name string, recordType RecordType) error {
	for i, record := range hz.Records {
		if record.Name == name && record.Type == recordType {
			hz.Records = append(hz.Records[:i], hz.Records[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("record not found: %s %s", name, recordType)
}

func (hz *HostedZone) CalculateCost(queriesPerMonth int64) float64 {
	zoneCost := hz.CostPerZone
	queryCost := float64(queriesPerMonth) * hz.CostPerQuery
	return zoneCost + queryCost
}

type RoutingPolicy string

const (
	RoutingPolicySimple           RoutingPolicy = "simple"
	RoutingPolicyWeighted         RoutingPolicy = "weighted"
	RoutingPolicyLatency          RoutingPolicy = "latency"
	RoutingPolicyFailover         RoutingPolicy = "failover"
	RoutingPolicyGeolocation      RoutingPolicy = "geolocation"
	RoutingPolicyGeoproximity     RoutingPolicy = "geoproximity"
	RoutingPolicyMultiValueAnswer RoutingPolicy = "multivalue"
)

type RoutingConfig struct {
	Policy        RoutingPolicy
	Weight        int
	Region        string
	SetID         string
	FailoverType  string
	HealthCheckID string
}

type LoadBalancerDNS struct {
	DNSName      string
	HostedZoneID string
	Region       string
	Type         string
}

func CreateLoadBalancerDNS(lbName, region string) *LoadBalancerDNS {
	return &LoadBalancerDNS{
		DNSName:      fmt.Sprintf("%s.%s.elb.amazonaws.com", lbName, region),
		HostedZoneID: fmt.Sprintf("Z%s", strings.ToUpper(region)),
		Region:       region,
		Type:         "application",
	}
}

type CDNDistribution struct {
	ID         string
	DomainName string
	Aliases    []string
	Origins    []*Origin
	Enabled    bool
	PriceClass string
	CostPerGB  float64
}

type Origin struct {
	ID                string
	DomainName        string
	Path              string
	CustomHeaders     map[string]string
	ConnectionTimeout int
	ResponseTimeout   int
}

func NewCDNDistribution(id string) *CDNDistribution {
	return &CDNDistribution{
		ID:         id,
		DomainName: fmt.Sprintf("%s.cloudfront.net", id),
		Aliases:    []string{},
		Origins:    []*Origin{},
		Enabled:    true,
		PriceClass: "PriceClass_All",
		CostPerGB:  0.085,
	}
}

func (cdn *CDNDistribution) AddOrigin(id, domainName, path string) *Origin {
	origin := &Origin{
		ID:                id,
		DomainName:        domainName,
		Path:              path,
		CustomHeaders:     make(map[string]string),
		ConnectionTimeout: 30,
		ResponseTimeout:   30,
	}
	cdn.Origins = append(cdn.Origins, origin)
	return origin
}

func (cdn *CDNDistribution) AddAlias(alias string) {
	cdn.Aliases = append(cdn.Aliases, alias)
}

func (cdn *CDNDistribution) CalculateCost(dataTransferGB float64) float64 {
	cost := dataTransferGB * cdn.CostPerGB

	if dataTransferGB > 10000 {
		cost = 10000*0.085 + (dataTransferGB-10000)*0.070
	} else if dataTransferGB > 40000 {
		cost = 10000*0.085 + 30000*0.070 + (dataTransferGB-40000)*0.050
	}

	return cost + 0.01
}

type DNSQueryLog struct {
	QueryName    string
	QueryType    RecordType
	Response     string
	ResponseTime time.Duration
	SourceIP     string
	Timestamp    time.Time
}

type DNSResolver struct {
	Cache        map[string]*CachedRecord
	MaxCacheSize int
}

type CachedRecord struct {
	Record    *DNSRecord
	CachedAt  time.Time
	ExpiresAt time.Time
}

func NewDNSResolver(maxCacheSize int) *DNSResolver {
	return &DNSResolver{
		Cache:        make(map[string]*CachedRecord),
		MaxCacheSize: maxCacheSize,
	}
}

func (resolver *DNSResolver) Resolve(name string, recordType RecordType, zone *HostedZone) (string, error) {
	cacheKey := fmt.Sprintf("%s:%s", name, recordType)

	if cached, exists := resolver.Cache[cacheKey]; exists {
		if time.Now().Before(cached.ExpiresAt) {
			return cached.Record.Value, nil
		}
		delete(resolver.Cache, cacheKey)
	}

	record := zone.GetRecord(name, recordType)
	if record == nil {
		return "", fmt.Errorf("record not found: %s %s", name, recordType)
	}

	if len(resolver.Cache) >= resolver.MaxCacheSize {
		resolver.evictOldest()
	}

	resolver.Cache[cacheKey] = &CachedRecord{
		Record:    record,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(record.TTL) * time.Second),
	}

	return record.Value, nil
}

func (resolver *DNSResolver) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, cached := range resolver.Cache {
		if oldestKey == "" || cached.CachedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.CachedAt
		}
	}

	if oldestKey != "" {
		delete(resolver.Cache, oldestKey)
	}
}

func (resolver *DNSResolver) ClearCache() {
	resolver.Cache = make(map[string]*CachedRecord)
}

func (resolver *DNSResolver) GetCacheSize() int {
	return len(resolver.Cache)
}

type DNSPreset struct {
	Name        string
	Description string
	Records     []PresetDNSRecord
}

type PresetDNSRecord struct {
	Name  string
	Type  RecordType
	Value string
	TTL   int
}

var DNSPresets = map[string]*DNSPreset{
	"simple-web": {
		Name:        "Simple Web Application",
		Description: "Basic DNS setup for a web application",
		Records: []PresetDNSRecord{
			{Name: "@", Type: RecordTypeA, Value: "192.0.2.1", TTL: 300},
			{Name: "www", Type: RecordTypeCNAME, Value: "@", TTL: 300},
			{Name: "@", Type: RecordTypeMX, Value: "10 mail.example.com", TTL: 3600},
		},
	},
	"cdn-enabled": {
		Name:        "CDN-Enabled Application",
		Description: "DNS setup with CDN",
		Records: []PresetDNSRecord{
			{Name: "@", Type: RecordTypeA, Value: "192.0.2.1", TTL: 300},
			{Name: "www", Type: RecordTypeCNAME, Value: "d111111abcdef8.cloudfront.net", TTL: 300},
			{Name: "static", Type: RecordTypeCNAME, Value: "d222222abcdef8.cloudfront.net", TTL: 300},
		},
	},
	"multi-region": {
		Name:        "Multi-Region Application",
		Description: "DNS setup with failover across regions",
		Records: []PresetDNSRecord{
			{Name: "us", Type: RecordTypeA, Value: "192.0.2.1", TTL: 60},
			{Name: "eu", Type: RecordTypeA, Value: "198.51.100.1", TTL: 60},
			{Name: "asia", Type: RecordTypeA, Value: "203.0.113.1", TTL: 60},
		},
	},
}

func GetDNSPreset(name string) *DNSPreset {
	if preset, exists := DNSPresets[name]; exists {
		return preset
	}
	return DNSPresets["simple-web"]
}

func GetDNSPresetNames() []string {
	return []string{"simple-web", "cdn-enabled", "multi-region"}
}

func CreateHostedZoneFromPreset(id, zoneName string, presetName string, vpc *VPC) *HostedZone {
	preset := GetDNSPreset(presetName)

	zone := NewHostedZone(id, zoneName, vpc != nil, vpc)

	for _, record := range preset.Records {
		recordName := record.Name
		if recordName == "@" {
			recordName = zoneName
		} else {
			recordName = fmt.Sprintf("%s.%s", record.Name, zoneName)
		}
		zone.AddRecord(recordName, record.Type, record.Value, record.TTL)
	}

	return zone
}

func ValidateDomainName(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	labels := strings.Split(domain, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}

		if !strings.ContainsAny(label[0:1], "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") {
			return false
		}
	}

	return true
}
