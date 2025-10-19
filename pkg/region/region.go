package region

type Region string

const (
	Global  Region = "global"
	America Region = "america"
	Europe  Region = "europe"
	Asia    Region = "asia"
	Oceania Region = "oceania"
)

func (r Region) String() string {
	return string(r)
}
