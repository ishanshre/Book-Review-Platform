package forms

// create a errors type that holds a map with index string and value a slice of strings
type errors map[string][]string

// Add is used to add an error to erros
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get is used to get the first error if exists
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
