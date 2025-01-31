package internal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/rivo/tview"
)

const (
	DB_DRIVER_FIELD = "DB driver"
	DB_URI_FIELD    = "DB Connection URI"
	DB_QUERY_FIELD  = "DB query"
	STATUS_FIELD    = "status"
	CONNECT_BUTTON  = "Connect Button"
	SAVE_BUTTON     = "Save Button"
	QUIT_BUTTON     = "Quit Button"
	PGSQLOPTION     = "PostgreSQL"
	MYSQLOPTION     = "MySQL"
	FIELD_WIDTH     = 100
)

type Console struct {
	app    *tview.Application
	form   *tview.Form
	driver *tview.DropDown
	dburi  *tview.InputField
	query  *tview.TextArea
	status *tview.TextView
	db     *sql.DB
	logbuf string
}

func NewConsole() *Console {
	app := tview.NewApplication()
	form := tview.NewForm()

	return &Console{
		app:  app,
		form: form,
	}

}

func (c *Console) OpenDB() {
	_, x := c.driver.GetCurrentOption()
	c.status.SetText(x)

	//hardcode for now
	host := "127.0.0.1"
	user := "graph"
	passwd := "graph"
	dbname := "graph"

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, passwd, host, "5432", dbname)
	db, err := NewDB("postgres", uri)
	if err != nil {
		c.status.SetText(fmt.Sprintf("Tried opening DB but got error: %s", err))
		return
	}

	if err := db.Ping(); err != nil {
		c.status.SetText(fmt.Sprintf("Tried to ping DB but got error %s", err))
		return
	}

	c.status.SetText(fmt.Sprintf("Connected %s", db))
}

func (c *Console) Run() {
	c.setLayout()
	if err := c.app.SetRoot(c.form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}

func (c *Console) Close() {
	log.Println("closing app...")
	c.app.Stop()
}

func (c *Console) Save() {
	log.Println("saving config...")
	c.app.Stop()
}

func (c *Console) GetValues() {
}

func (c *Console) Connect() {
	c.OpenDB()

}

func (c *Console) addstatus(logline string) {

}

func (c *Console) changeDriver(label string, index int) {
	if label == PGSQLOPTION {
		c.dburi.SetText("postgres://<username>:<password>@<host>/<dbname>?sslmode=<verify,disable>")
	}

	if label == MYSQLOPTION {
		c.dburi.SetText("<username>:<password>@<host:port>/<dbname>?<paramN=valueN,...>")
	}
	c.app.SetFocus(c.dburi)
}

func (c *Console) setLayout() {

	c.form.AddInputField(DB_URI_FIELD, "", FIELD_WIDTH, nil, nil)
	c.dburi = c.form.GetFormItemByLabel(DB_URI_FIELD).(*tview.InputField)

	c.form.AddDropDown(DB_DRIVER_FIELD, []string{PGSQLOPTION, MYSQLOPTION}, 0, c.changeDriver)
	c.driver = c.form.GetFormItemByLabel(DB_DRIVER_FIELD).(*tview.DropDown)

	c.form.AddTextArea(DB_QUERY_FIELD, "", FIELD_WIDTH, 10, 500, nil)
	c.query = c.form.GetFormItemByLabel(DB_QUERY_FIELD).(*tview.TextArea)

	c.form.AddTextView(STATUS_FIELD, "", FIELD_WIDTH, 10, false, true)
	c.status = c.form.GetFormItemByLabel(STATUS_FIELD).(*tview.TextView)

	c.form.AddButton(CONNECT_BUTTON, c.Connect)
	c.form.AddButton(SAVE_BUTTON, c.Save)
	c.form.AddButton(QUIT_BUTTON, c.Close)

}
