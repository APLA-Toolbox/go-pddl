package common

type StringSlice []string

func (ss StringSlice) ToString(name string) string {
	s := name + ": ["
	for _, v := range ss {
		s += "'" + v + "'"
	}
	s += "]\n"
	return s
}

type StringMap map[string]string