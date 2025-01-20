package internal

import (
	"log"

	"github.com/rivo/tview"
)

type Console struct {
	App    *tview.Application
	Form   *tview.Form
	Host   *tview.InputField
	User   *tview.InputField
	Passwd *tview.InputField
	DBName *tview.InputField
	Tables *tview.DropDown
	Status *tview.InputField
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
	status := c.Form.GetFormItemByLabel("Status").(*tview.InputField)
	status.SetText("Connecting To Database ...")
	log.Println(pass.GetText())

}

func (c *Console) SetLayout() {
	c.Form.AddDropDown("DB System", []string{"PostgreSQL", "MySQL"}, 0, nil)
	c.Form.AddInputField("DB Host", "", 50, nil, nil)
	c.Form.AddInputField("DB User", "", 50, nil, nil)
	c.Form.AddPasswordField("DB Password", "", 50, '*', nil)
	c.Form.AddInputField("DB Name", "", 50, nil, nil)
	c.Form.AddDropDown("DB Tables", []string{"Connect to populate"}, 0, nil)
	c.Form.AddButton("Connect To DB", c.Connect)
	c.Form.AddInputField("Status", "", 50, nil, nil)
	c.Form.AddButton("Save Configuration", c.Save)
	c.Form.AddButton("Quit", c.Close)

	//save references to tview elements so we dont have to keep looking up
	c.Host = c.Form.GetFormItemByLabel("DB Host").(*tview.InputField)
	c.User = c.Form.GetFormItemByLabel("DB User").(*tview.InputField)
	c.Passwd = c.Form.GetFormItemByLabel("DB Password").(*tview.InputField)
	c.DBName = c.Form.GetFormItemByLabel("DB Name").(*tview.InputField)
	c.Tables = c.Form.GetFormItemByLabel("DB Tables").(*tview.DropDown)
	c.Status = c.Form.GetFormItemByLabel("Status").(*tview.InputField)

	//status input has to be disabled as its only for print status
	c.Status.SetDisabled(true)

	//set focus to first input
	c.App.SetFocus(c.Form.GetFormItemByLabel("DB Host"))
}
