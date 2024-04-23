package core

type (
	StringPool map[string]*string
)

func (p StringPool) Intern(s string) string {
	return s
	/*
		ss, exists := p[s]
		if exists {
			return ss
		}
		p[s] = &s
		return &s
	*/
}
