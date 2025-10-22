package region

type Region string

const (
	Global  Region = "global"
	Africa  Region = "africa"
	America Region = "america"
	Asia    Region = "asia"
	Europe  Region = "europe"
	Oceania Region = "oceania"
)

const (
	bitAfrica  = 1 << iota // 00001
	bitAmerica             // 00010
	bitAsia                // 00100
	bitEurope              // 01000
	bitOceania             // 10000
)

var regionBitMap = map[Region]int{
	Africa:  bitAfrica,
	America: bitAmerica,
	Asia:    bitAsia,
	Europe:  bitEurope,
	Oceania: bitOceania,
}

func FromInt(mask int) []Region {
	var regions []Region
	for region, bit := range regionBitMap {
		if mask&bit != 0 {
			regions = append(regions, region)
		}
	}
	return regions
}

func ToInt(regions []Region) int {
	var mask int
	for _, region := range regions {
		if bit, ok := regionBitMap[region]; ok {
			mask |= bit
		}
	}
	return mask
}
