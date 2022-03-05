package refactor

// Refactor entrypoint
func DoRefactor(input string) string {
	if input == "foo" {
		return "refactor bar"
	}
	return "refactor nothing"
}
