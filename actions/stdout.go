package actions

import (
	"fmt"
	"os"

	"github.com/khulnasoft-lab/hooker/v2/formatting"
	"github.com/khulnasoft-lab/hooker/v2/layout"
)

type StdoutAction struct {
	Name string
}

func (stdout StdoutAction) GetName() string { return stdout.Name }
func (stdout StdoutAction) Init() error {
	return nil
}
func (stdout StdoutAction) Send(data map[string]string) error {
	_, err := fmt.Fprintf(os.Stdout, "%s", data["description"])
	return err
}
func (stdout StdoutAction) Terminate() error {
	return nil
}
func (stdout StdoutAction) GetLayoutProvider() layout.LayoutProvider {
	return &formatting.HtmlProvider{}
}
