package config

import "time"

type Region struct {
	ID                string
	Name              string
	DisplayName       string
	Location          string
	Latitude          float64
	Longitude         float64
	AvailabilityZones []string
	Active            bool
}

var Regions = map[string]*Region{
	"us-east-1": {
		ID:                "us-east-1",
		Name:              "us-east",
		DisplayName:       "US East (N. Virginia)",
		Location:          "Virginia, USA",
		Latitude:          38.13,
		Longitude:         -78.45,
		AvailabilityZones: []string{"us-east-1a", "us-east-1b", "us-east-1c", "us-east-1d"},
		Active:            true,
	},
	"us-west-1": {
		ID:                "us-west-1",
		Name:              "us-west",
		DisplayName:       "US West (N. California)",
		Location:          "California, USA",
		Latitude:          37.35,
		Longitude:         -121.96,
		AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
		Active:            true,
	},
	"eu-west-1": {
		ID:                "eu-west-1",
		Name:              "europe",
		DisplayName:       "EU West (Ireland)",
		Location:          "Dublin, Ireland",
		Latitude:          53.35,
		Longitude:         -6.26,
		AvailabilityZones: []string{"eu-west-1a", "eu-west-1b", "eu-west-1c"},
		Active:            true,
	},
	"ap-southeast-1": {
		ID:                "ap-southeast-1",
		Name:              "asia",
		DisplayName:       "Asia Pacific (Singapore)",
		Location:          "Singapore",
		Latitude:          1.29,
		Longitude:         103.85,
		AvailabilityZones: []string{"ap-southeast-1a", "ap-southeast-1b", "ap-southeast-1c"},
		Active:            true,
	},
	"ap-southeast-2": {
		ID:                "ap-southeast-2",
		Name:              "australia",
		DisplayName:       "Asia Pacific (Sydney)",
		Location:          "Sydney, Australia",
		Latitude:          -33.86,
		Longitude:         151.20,
		AvailabilityZones: []string{"ap-southeast-2a", "ap-southeast-2b", "ap-southeast-2c"},
		Active:            true,
	},
}

func GetRegion(id string) *Region {
	if region, exists := Regions[id]; exists {
		return region
	}
	return Regions["us-east-1"]
}

func GetRegionByName(name string) *Region {
	for _, region := range Regions {
		if region.Name == name {
			return region
		}
	}
	return Regions["us-east-1"]
}

func GetAllRegions() []*Region {
	regions := []*Region{}
	for _, region := range Regions {
		if region.Active {
			regions = append(regions, region)
		}
	}
	return regions
}

func GetRegionNames() []string {
	names := []string{}
	for _, region := range Regions {
		if region.Active {
			names = append(names, region.Name)
		}
	}
	return names
}

func GetRegionIDs() []string {
	ids := []string{"us-east-1", "us-west-1", "eu-west-1", "ap-southeast-1", "ap-southeast-2"}
	return ids
}

type NetworkLatency struct {
	FromRegion string
	ToRegion   string
	Latency    time.Duration
}

var NetworkLatencyMap = map[string]time.Duration{
	"us-east-1:us-east-1":      5 * time.Millisecond,
	"us-east-1:us-west-1":      70 * time.Millisecond,
	"us-east-1:eu-west-1":      90 * time.Millisecond,
	"us-east-1:ap-southeast-1": 150 * time.Millisecond,
	"us-east-1:ap-southeast-2": 180 * time.Millisecond,

	"us-west-1:us-west-1":      5 * time.Millisecond,
	"us-west-1:us-east-1":      70 * time.Millisecond,
	"us-west-1:eu-west-1":      140 * time.Millisecond,
	"us-west-1:ap-southeast-1": 120 * time.Millisecond,
	"us-west-1:ap-southeast-2": 130 * time.Millisecond,

	"eu-west-1:eu-west-1":      5 * time.Millisecond,
	"eu-west-1:us-east-1":      90 * time.Millisecond,
	"eu-west-1:us-west-1":      140 * time.Millisecond,
	"eu-west-1:ap-southeast-1": 120 * time.Millisecond,
	"eu-west-1:ap-southeast-2": 220 * time.Millisecond,

	"ap-southeast-1:ap-southeast-1": 5 * time.Millisecond,
	"ap-southeast-1:us-east-1":      150 * time.Millisecond,
	"ap-southeast-1:us-west-1":      120 * time.Millisecond,
	"ap-southeast-1:eu-west-1":      120 * time.Millisecond,
	"ap-southeast-1:ap-southeast-2": 90 * time.Millisecond,

	"ap-southeast-2:ap-southeast-2": 5 * time.Millisecond,
	"ap-southeast-2:us-east-1":      180 * time.Millisecond,
	"ap-southeast-2:us-west-1":      130 * time.Millisecond,
	"ap-southeast-2:eu-west-1":      220 * time.Millisecond,
	"ap-southeast-2:ap-southeast-1": 90 * time.Millisecond,
}

func GetNetworkLatency(fromRegion, toRegion string) time.Duration {
	key := fromRegion + ":" + toRegion
	if latency, exists := NetworkLatencyMap[key]; exists {
		return latency
	}

	reverseKey := toRegion + ":" + fromRegion
	if latency, exists := NetworkLatencyMap[reverseKey]; exists {
		return latency
	}

	return 100 * time.Millisecond
}

type AvailabilityZone struct {
	ID     string
	Region string
	Name   string
	Letter string
}

func GetAvailabilityZones(regionID string) []string {
	if region, exists := Regions[regionID]; exists {
		return region.AvailabilityZones
	}
	return []string{}
}

func GetAvailabilityZone(azID string) *AvailabilityZone {
	for regionID, region := range Regions {
		for _, az := range region.AvailabilityZones {
			if az == azID {
				letter := string(azID[len(azID)-1])
				return &AvailabilityZone{
					ID:     azID,
					Region: regionID,
					Name:   region.DisplayName + " - AZ " + letter,
					Letter: letter,
				}
			}
		}
	}
	return nil
}

type DataCenter struct {
	Region           string
	AvailabilityZone string
	Name             string
}

func NewDataCenter(region, az string) *DataCenter {
	return &DataCenter{
		Region:           region,
		AvailabilityZone: az,
		Name:             region + "/" + az,
	}
}
