package utils

type Strings []string

func (slice Strings) Include(value string) bool {
	for _, item := range slice {
		if value == item {
			return true
		}
	}
	return false
}

func (slice Strings) Unique() []string {
	unique := Strings{}
	for _, v := range slice {
		if !unique.Include(v) {
			unique = append(unique, v)
		}
	}

	return unique
}
