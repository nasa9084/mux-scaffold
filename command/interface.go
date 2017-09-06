package command

// Command interface
type Command interface {
	Exec(args []string) int
	Description() string
	Help()
}
