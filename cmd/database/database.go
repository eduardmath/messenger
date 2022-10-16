package database

type Database struct {
	s string
}

func (db Database) Print(s string) string {
	var ss string
	var len = len(s)
	for i := 0; i < len; i++ {
		ss += string(s[len-i-1])
	}

	return ss
}
