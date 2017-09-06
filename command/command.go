package command

// List of available commands
var List = map[string]Command{
	"help": &Help{},
	"init": &Init{},
}
