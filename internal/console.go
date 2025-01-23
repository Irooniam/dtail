package internal

import (
	"log"

	"github.com/rivo/tview"
)

const (
	DB_SYSTEM_FIELD = "DB System"
	DB_HOST_FIELD   = "DB Host"
	DB_USER_FIELD   = "DB User"
	DB_PASSWD_FIELD = "DB Password"
	DB_NAME_FIELD   = "DB Name"
	DB_QUERY_FIELD  = "DB Query"
	STATUS_FIELD    = "Status"
	CONNECT_BUTTON  = "Connect Button"
	SAVE_BUTTON     = "Save Button"
	QUIT_BUTTON     = "Quit Button"
)

type Console struct {
	App    *tview.Application
	Form   *tview.Form
	System *tview.DropDown
	Host   *tview.InputField
	User   *tview.InputField
	Passwd *tview.InputField
	DBName *tview.InputField
	Query  *tview.TextArea
	Status *tview.TextArea
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

func (c *Console) GetValues() {
}

func (c *Console) Connect() {
	system := c.System.GetTitle()
	log.Println(system)

}

func (c *Console) SetLayout() {
	c.Form.AddDropDown(DB_SYSTEM_FIELD, []string{"PostgreSQL", "MySQL"}, 0, nil)
	c.Form.AddInputField(DB_HOST_FIELD, "", 50, nil, nil)
	c.Form.AddInputField(DB_USER_FIELD, "", 50, nil, nil)
	c.Form.AddPasswordField(DB_PASSWD_FIELD, "", 50, '*', nil)
	c.Form.AddInputField(DB_NAME_FIELD, "", 50, nil, nil)
	c.Form.AddTextArea(DB_QUERY_FIELD, "", 50, 10, 500, nil)
	c.Form.AddTextArea(STATUS_FIELD, "", 50, 10, 500, nil)

	c.Form.AddButton(CONNECT_BUTTON, c.Connect)
	c.Form.AddButton(SAVE_BUTTON, c.Save)
	c.Form.AddButton(QUIT_BUTTON, c.Close)

	//save references to tview elements so we dont have to keep looking up
	c.System = c.Form.GetFormItemByLabel(DB_SYSTEM_FIELD).(*tview.DropDown)
	c.Host = c.Form.GetFormItemByLabel(DB_HOST_FIELD).(*tview.InputField)
	c.User = c.Form.GetFormItemByLabel(DB_USER_FIELD).(*tview.InputField)
	c.Passwd = c.Form.GetFormItemByLabel(DB_PASSWD_FIELD).(*tview.InputField)
	c.DBName = c.Form.GetFormItemByLabel(DB_NAME_FIELD).(*tview.InputField)
	c.Query = c.Form.GetFormItemByLabel(DB_QUERY_FIELD).(*tview.TextArea)
	c.Status = c.Form.GetFormItemByLabel(STATUS_FIELD).(*tview.TextArea)

	//status input has to be disabled as its only for print status
	c.Status.SetDisabled(true)

	//set focus to first input
	c.App.SetFocus(c.Form.GetFormItemByLabel("DB Host"))
}
