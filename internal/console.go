package internal

import (
	"log"

	"github.com/rivo/tview"
)

type Console struct {
	App     *tview.Application
	Form    *tview.Form
	Header  string
	FooterB string
	HostB   string
	UserB   string
	PassB   string
	DBNameB string
	TablesB string
}

func NewConsole() *Console {
	log.Println("new console...")

	app := tview.NewApplication()
	form := tview.NewForm()
	//app.SetRoot(form, true).EnableMouse(true).EnablePaste(true)
	//app.SetRoot(form, true).EnablePaste(true)

	return &Console{
		App:  app,
		Form: form,
	}

}

func (c *Console) Close() {
	log.Println("closing app...")
	c.App.Stop()
}

func (c *Console) Save() {
	log.Println("saving config...")
	c.App.Stop()
}

func (c *Console) Connect() {
	log.Println("connecting...")
	pass := c.Form.GetFormItemByLabel("DB Password").(*tview.InputField)
	log.Println(pass.GetText())

}

func (c *Console) SetLayout() {
	c.Form.AddInputField("DB Host", "", 50, nil, nil)
	c.Form.AddInputField("DB User", "", 50, nil, nil)
	c.Form.AddPasswordField("DB Password", "", 50, '*', nil)
	c.Form.AddInputField("DB Name", "", 50, nil, nil)
	c.Form.AddDropDown("DB Tables", []string{"Connect to populate"}, 0, nil)
	c.Form.AddButton("Connect", c.Connect)
	c.Form.AddButton("Save", c.Save)
	c.Form.AddButton("Quit", c.Close)

	c.App.SetFocus(c.Form.GetFormItemByLabel("DB Host"))
}
