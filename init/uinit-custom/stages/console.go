package stages

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"uinit-custom/config"
)

//Console implementes IStage
type Console struct {
	mode   string
	prefix string
	c      config.Config
	s      config.Secrets
	finals []string
}

//String ..
func (con Console) String() string {
	return "Console"
}

//Finalise ..
func (con Console) Finalise() []string {
	return con.finals
}

//Run ..
func (con Console) Run(c config.Config, s config.Secrets) error {

	con.c = c
	con.s = s
	con.mode = "ENG"
	con.prefix = "mjolnir"

	commands := []string{}

	reader := bufio.NewReader(os.Stdin)
	flag := ""
	con.write("Enter ? for commands")

	for flag != "q\n" {

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("mjolnir: failed to read string: %v", err)
		}
		flag = text

		response := con.parse(text)
		con.write(response)

	}

	err := execute(commands)
	if err != nil {
		return err
	}

	return nil
}

func (con Console) write(str string) {

	fmt.Printf("[%v] %v: %v\n", con.mode, con.prefix, str)
	fmt.Printf("[%v] %v> ", con.mode, con.prefix)
}

func (con Console) parse(command string) string {

	response := "Invalid command, enter ? for commands"

	switch con.mode {
	case "ENG":
		switch command {
		case "?\n":
			{
				return fmt.Sprintf("\n?: help\nq: quit\ndc: dump config\nsn: show network\nsw: show wireguard")
			}
		case "sn\n":
			{
				str, err := executeOne("ip a", "")
				if err != nil {
					return fmt.Sprintf("%v: %v", str, err)
				}
				return str
			}
		case "sw\n":
			{
				str, err := executeOne("wg show", "")
				if err != nil {
					return fmt.Sprintf("%v: %v", str, err)
				}
				return str
			}
		case "dc\n":
			{
				json, err := json.MarshalIndent(con.c, "", "    ")
				if err != nil {
					return fmt.Sprintf("Failed to marshal config: %v", err)
				}
				return string(json)
			}
		default:
			{
				return response
			}
		}
	}
	return "invalid selection"
}
