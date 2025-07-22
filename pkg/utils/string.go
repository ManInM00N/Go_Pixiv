package utils

var pattern = [...]string{"R-18", "r-18", "r18", "R18"}

func HasR18(raw *[]string) bool {
	for _, v := range *raw {
		if len(v) > 4 || len(v) < 3 {
			continue
		}
		for _, p := range pattern {
			if v == p {
				return true
			}
		}
	}
	return false
}
