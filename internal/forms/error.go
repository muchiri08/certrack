package forms

type errros map[string][]string

func (e errros) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errros) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
