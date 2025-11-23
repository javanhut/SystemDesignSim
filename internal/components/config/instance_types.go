package config

type InstanceType struct {
	Name              string
	DisplayName       string
	Category          InstanceCategory
	vCPU              int
	MemoryGB          float64
	NetworkGbps       float64
	StorageGB         int
	CostPerHour       float64
	MaxRequestsPerSec int
	Description       string
}

type InstanceCategory string

const (
	CategoryBurstable      InstanceCategory = "burstable"
	CategoryGeneralPurpose InstanceCategory = "general-purpose"
	CategoryComputeOpt     InstanceCategory = "compute-optimized"
	CategoryMemoryOpt      InstanceCategory = "memory-optimized"
	CategoryStorageOpt     InstanceCategory = "storage-optimized"
)

var InstanceTypes = map[string]*InstanceType{
	"t2.micro": {
		Name:              "t2.micro",
		DisplayName:       "t2.micro (1 vCPU, 1GB RAM)",
		Category:          CategoryBurstable,
		vCPU:              1,
		MemoryGB:          1.0,
		NetworkGbps:       0.1,
		StorageGB:         8,
		CostPerHour:       0.0116,
		MaxRequestsPerSec: 10,
		Description:       "Burstable - Good for low traffic applications",
	},
	"t2.small": {
		Name:              "t2.small",
		DisplayName:       "t2.small (1 vCPU, 2GB RAM)",
		Category:          CategoryBurstable,
		vCPU:              1,
		MemoryGB:          2.0,
		NetworkGbps:       0.2,
		StorageGB:         8,
		CostPerHour:       0.023,
		MaxRequestsPerSec: 20,
		Description:       "Burstable - Good for small applications",
	},
	"t2.medium": {
		Name:              "t2.medium",
		DisplayName:       "t2.medium (2 vCPU, 4GB RAM)",
		Category:          CategoryBurstable,
		vCPU:              2,
		MemoryGB:          4.0,
		NetworkGbps:       0.3,
		StorageGB:         16,
		CostPerHour:       0.0464,
		MaxRequestsPerSec: 40,
		Description:       "Burstable - Good for moderate traffic",
	},
	"t3.micro": {
		Name:              "t3.micro",
		DisplayName:       "t3.micro (2 vCPU, 1GB RAM)",
		Category:          CategoryBurstable,
		vCPU:              2,
		MemoryGB:          1.0,
		NetworkGbps:       0.5,
		StorageGB:         8,
		CostPerHour:       0.0104,
		MaxRequestsPerSec: 15,
		Description:       "Burstable - Latest generation, better performance",
	},
	"t3.small": {
		Name:              "t3.small",
		DisplayName:       "t3.small (2 vCPU, 2GB RAM)",
		Category:          CategoryBurstable,
		vCPU:              2,
		MemoryGB:          2.0,
		NetworkGbps:       0.5,
		StorageGB:         8,
		CostPerHour:       0.0208,
		MaxRequestsPerSec: 30,
		Description:       "Burstable - Latest generation",
	},
	"t3.medium": {
		Name:              "t3.medium",
		DisplayName:       "t3.medium (2 vCPU, 4GB RAM)",
		Category:          CategoryBurstable,
		vCPU:              2,
		MemoryGB:          4.0,
		NetworkGbps:       0.5,
		StorageGB:         16,
		CostPerHour:       0.0416,
		MaxRequestsPerSec: 50,
		Description:       "Burstable - Latest generation",
	},
	"m5.large": {
		Name:              "m5.large",
		DisplayName:       "m5.large (2 vCPU, 8GB RAM)",
		Category:          CategoryGeneralPurpose,
		vCPU:              2,
		MemoryGB:          8.0,
		NetworkGbps:       1.0,
		StorageGB:         32,
		CostPerHour:       0.096,
		MaxRequestsPerSec: 100,
		Description:       "General Purpose - Balanced compute, memory, network",
	},
	"m5.xlarge": {
		Name:              "m5.xlarge",
		DisplayName:       "m5.xlarge (4 vCPU, 16GB RAM)",
		Category:          CategoryGeneralPurpose,
		vCPU:              4,
		MemoryGB:          16.0,
		NetworkGbps:       1.25,
		StorageGB:         64,
		CostPerHour:       0.192,
		MaxRequestsPerSec: 200,
		Description:       "General Purpose - Good for most workloads",
	},
	"m5.2xlarge": {
		Name:              "m5.2xlarge",
		DisplayName:       "m5.2xlarge (8 vCPU, 32GB RAM)",
		Category:          CategoryGeneralPurpose,
		vCPU:              8,
		MemoryGB:          32.0,
		NetworkGbps:       2.5,
		StorageGB:         128,
		CostPerHour:       0.384,
		MaxRequestsPerSec: 400,
		Description:       "General Purpose - High performance",
	},
	"m5.4xlarge": {
		Name:              "m5.4xlarge",
		DisplayName:       "m5.4xlarge (16 vCPU, 64GB RAM)",
		Category:          CategoryGeneralPurpose,
		vCPU:              16,
		MemoryGB:          64.0,
		NetworkGbps:       5.0,
		StorageGB:         256,
		CostPerHour:       0.768,
		MaxRequestsPerSec: 800,
		Description:       "General Purpose - Very high performance",
	},
	"c5.large": {
		Name:              "c5.large",
		DisplayName:       "c5.large (2 vCPU, 4GB RAM)",
		Category:          CategoryComputeOpt,
		vCPU:              2,
		MemoryGB:          4.0,
		NetworkGbps:       1.0,
		StorageGB:         16,
		CostPerHour:       0.085,
		MaxRequestsPerSec: 120,
		Description:       "Compute Optimized - CPU intensive workloads",
	},
	"c5.xlarge": {
		Name:              "c5.xlarge",
		DisplayName:       "c5.xlarge (4 vCPU, 8GB RAM)",
		Category:          CategoryComputeOpt,
		vCPU:              4,
		MemoryGB:          8.0,
		NetworkGbps:       1.25,
		StorageGB:         32,
		CostPerHour:       0.17,
		MaxRequestsPerSec: 240,
		Description:       "Compute Optimized - High CPU performance",
	},
	"c5.2xlarge": {
		Name:              "c5.2xlarge",
		DisplayName:       "c5.2xlarge (8 vCPU, 16GB RAM)",
		Category:          CategoryComputeOpt,
		vCPU:              8,
		MemoryGB:          16.0,
		NetworkGbps:       2.5,
		StorageGB:         64,
		CostPerHour:       0.34,
		MaxRequestsPerSec: 480,
		Description:       "Compute Optimized - Very high CPU",
	},
	"r5.large": {
		Name:              "r5.large",
		DisplayName:       "r5.large (2 vCPU, 16GB RAM)",
		Category:          CategoryMemoryOpt,
		vCPU:              2,
		MemoryGB:          16.0,
		NetworkGbps:       1.0,
		StorageGB:         32,
		CostPerHour:       0.126,
		MaxRequestsPerSec: 90,
		Description:       "Memory Optimized - Memory intensive workloads",
	},
	"r5.xlarge": {
		Name:              "r5.xlarge",
		DisplayName:       "r5.xlarge (4 vCPU, 32GB RAM)",
		Category:          CategoryMemoryOpt,
		vCPU:              4,
		MemoryGB:          32.0,
		NetworkGbps:       1.25,
		StorageGB:         64,
		CostPerHour:       0.252,
		MaxRequestsPerSec: 180,
		Description:       "Memory Optimized - High memory applications",
	},
	"r5.2xlarge": {
		Name:              "r5.2xlarge",
		DisplayName:       "r5.2xlarge (8 vCPU, 64GB RAM)",
		Category:          CategoryMemoryOpt,
		vCPU:              8,
		MemoryGB:          64.0,
		NetworkGbps:       2.5,
		StorageGB:         128,
		CostPerHour:       0.504,
		MaxRequestsPerSec: 360,
		Description:       "Memory Optimized - Very high memory",
	},
}

func GetInstanceType(name string) *InstanceType {
	if it, exists := InstanceTypes[name]; exists {
		return it
	}
	return InstanceTypes["t2.micro"]
}

func GetInstanceTypesByCategory(category InstanceCategory) []*InstanceType {
	types := []*InstanceType{}
	for _, it := range InstanceTypes {
		if it.Category == category {
			types = append(types, it)
		}
	}
	return types
}

func GetAllInstanceTypes() []*InstanceType {
	types := []*InstanceType{}
	for _, it := range InstanceTypes {
		types = append(types, it)
	}
	return types
}

func GetInstanceTypeNames() []string {
	names := []string{
		"t2.micro", "t2.small", "t2.medium",
		"t3.micro", "t3.small", "t3.medium",
		"m5.large", "m5.xlarge", "m5.2xlarge", "m5.4xlarge",
		"c5.large", "c5.xlarge", "c5.2xlarge",
		"r5.large", "r5.xlarge", "r5.2xlarge",
	}
	return names
}

type DatabaseInstanceType struct {
	Name             string
	DisplayName      string
	vCPU             int
	MemoryGB         float64
	StorageGB        int
	IOPS             int
	MaxConnections   int
	CostPerHour      float64
	CostPerGBStorage float64
	Description      string
}

var DatabaseInstanceTypes = map[string]*DatabaseInstanceType{
	"db.t2.micro": {
		Name:             "db.t2.micro",
		DisplayName:      "db.t2.micro (1 vCPU, 1GB RAM)",
		vCPU:             1,
		MemoryGB:         1.0,
		StorageGB:        20,
		IOPS:             100,
		MaxConnections:   50,
		CostPerHour:      0.017,
		CostPerGBStorage: 0.10,
		Description:      "Small database - Development/testing",
	},
	"db.t2.small": {
		Name:             "db.t2.small",
		DisplayName:      "db.t2.small (1 vCPU, 2GB RAM)",
		vCPU:             1,
		MemoryGB:         2.0,
		StorageGB:        20,
		IOPS:             200,
		MaxConnections:   100,
		CostPerHour:      0.034,
		CostPerGBStorage: 0.10,
		Description:      "Small database - Low traffic apps",
	},
	"db.t3.small": {
		Name:             "db.t3.small",
		DisplayName:      "db.t3.small (2 vCPU, 2GB RAM)",
		vCPU:             2,
		MemoryGB:         2.0,
		StorageGB:        20,
		IOPS:             250,
		MaxConnections:   150,
		CostPerHour:      0.034,
		CostPerGBStorage: 0.10,
		Description:      "Latest gen - Better performance",
	},
	"db.t3.medium": {
		Name:             "db.t3.medium",
		DisplayName:      "db.t3.medium (2 vCPU, 4GB RAM)",
		vCPU:             2,
		MemoryGB:         4.0,
		StorageGB:        50,
		IOPS:             500,
		MaxConnections:   300,
		CostPerHour:      0.068,
		CostPerGBStorage: 0.10,
		Description:      "Medium database - Moderate traffic",
	},
	"db.m5.large": {
		Name:             "db.m5.large",
		DisplayName:      "db.m5.large (2 vCPU, 8GB RAM)",
		vCPU:             2,
		MemoryGB:         8.0,
		StorageGB:        100,
		IOPS:             1000,
		MaxConnections:   500,
		CostPerHour:      0.186,
		CostPerGBStorage: 0.115,
		Description:      "Production database - Good performance",
	},
	"db.m5.xlarge": {
		Name:             "db.m5.xlarge",
		DisplayName:      "db.m5.xlarge (4 vCPU, 16GB RAM)",
		vCPU:             4,
		MemoryGB:         16.0,
		StorageGB:        200,
		IOPS:             2000,
		MaxConnections:   1000,
		CostPerHour:      0.372,
		CostPerGBStorage: 0.115,
		Description:      "Production database - High performance",
	},
	"db.m5.2xlarge": {
		Name:             "db.m5.2xlarge",
		DisplayName:      "db.m5.2xlarge (8 vCPU, 32GB RAM)",
		vCPU:             8,
		MemoryGB:         32.0,
		StorageGB:        500,
		IOPS:             5000,
		MaxConnections:   2000,
		CostPerHour:      0.744,
		CostPerGBStorage: 0.115,
		Description:      "Production database - Very high performance",
	},
	"db.r5.large": {
		Name:             "db.r5.large",
		DisplayName:      "db.r5.large (2 vCPU, 16GB RAM)",
		vCPU:             2,
		MemoryGB:         16.0,
		StorageGB:        100,
		IOPS:             1000,
		MaxConnections:   700,
		CostPerHour:      0.24,
		CostPerGBStorage: 0.115,
		Description:      "Memory optimized - Large datasets",
	},
	"db.r5.xlarge": {
		Name:             "db.r5.xlarge",
		DisplayName:      "db.r5.xlarge (4 vCPU, 32GB RAM)",
		vCPU:             4,
		MemoryGB:         32.0,
		StorageGB:        200,
		IOPS:             2000,
		MaxConnections:   1500,
		CostPerHour:      0.48,
		CostPerGBStorage: 0.115,
		Description:      "Memory optimized - Very large datasets",
	},
}

func GetDatabaseInstanceType(name string) *DatabaseInstanceType {
	if it, exists := DatabaseInstanceTypes[name]; exists {
		return it
	}
	return DatabaseInstanceTypes["db.t2.micro"]
}

func GetDatabaseInstanceTypeNames() []string {
	names := []string{
		"db.t2.micro", "db.t2.small",
		"db.t3.small", "db.t3.medium",
		"db.m5.large", "db.m5.xlarge", "db.m5.2xlarge",
		"db.r5.large", "db.r5.xlarge",
	}
	return names
}

type CacheInstanceType struct {
	Name           string
	DisplayName    string
	MemoryGB       float64
	NetworkGbps    float64
	MaxConnections int
	CostPerHour    float64
	Description    string
}

var CacheInstanceTypes = map[string]*CacheInstanceType{
	"cache.t2.micro": {
		Name:           "cache.t2.micro",
		DisplayName:    "cache.t2.micro (0.5GB)",
		MemoryGB:       0.5,
		NetworkGbps:    0.1,
		MaxConnections: 100,
		CostPerHour:    0.017,
		Description:    "Small cache - Development",
	},
	"cache.t3.small": {
		Name:           "cache.t3.small",
		DisplayName:    "cache.t3.small (1.5GB)",
		MemoryGB:       1.5,
		NetworkGbps:    0.5,
		MaxConnections: 300,
		CostPerHour:    0.034,
		Description:    "Small cache - Low traffic",
	},
	"cache.m5.large": {
		Name:           "cache.m5.large",
		DisplayName:    "cache.m5.large (6.4GB)",
		MemoryGB:       6.4,
		NetworkGbps:    1.0,
		MaxConnections: 1000,
		CostPerHour:    0.126,
		Description:    "Medium cache - Production",
	},
	"cache.m5.xlarge": {
		Name:           "cache.m5.xlarge",
		DisplayName:    "cache.m5.xlarge (12.9GB)",
		MemoryGB:       12.9,
		NetworkGbps:    1.25,
		MaxConnections: 2000,
		CostPerHour:    0.252,
		Description:    "Large cache - High traffic",
	},
	"cache.r5.large": {
		Name:           "cache.r5.large",
		DisplayName:    "cache.r5.large (13.1GB)",
		MemoryGB:       13.1,
		NetworkGbps:    1.0,
		MaxConnections: 1500,
		CostPerHour:    0.188,
		Description:    "Memory optimized - Large cache",
	},
	"cache.r5.xlarge": {
		Name:           "cache.r5.xlarge",
		DisplayName:    "cache.r5.xlarge (26.3GB)",
		MemoryGB:       26.3,
		NetworkGbps:    1.25,
		MaxConnections: 3000,
		CostPerHour:    0.376,
		Description:    "Memory optimized - Very large cache",
	},
}

func GetCacheInstanceType(name string) *CacheInstanceType {
	if it, exists := CacheInstanceTypes[name]; exists {
		return it
	}
	return CacheInstanceTypes["cache.t2.micro"]
}

func GetCacheInstanceTypeNames() []string {
	names := []string{
		"cache.t2.micro", "cache.t3.small",
		"cache.m5.large", "cache.m5.xlarge",
		"cache.r5.large", "cache.r5.xlarge",
	}
	return names
}
