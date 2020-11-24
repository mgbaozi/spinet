package main
import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "spinet-cli"
	app.Usage = "Spinet command line tools"
	app.Commands = []cli.Command{
	}
	app.Action = core
	app.Version = "0.0.1"
	app.Run(os.Args)
}
