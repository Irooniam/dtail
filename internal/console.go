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
	CONNECT_BUTTON  = "Connect DB"
	TABLES_BUTTON   = "List Tables"
	SAVE_BUTTON     = "Save"
	QUIT_BUTTON     = "Quit"
	PGSQLOPTION     = "postgres"
	MYSQLOPTION     = "mysql"
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

type dbInter interface {
	getTables() ([]string, error)
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
	_, driver := c.driver.GetCurrentOption()
	c.addStatus(fmt.Sprintf("Preparing to connect to database server using %s driver", driver))

	db, err := NewDB(driver, c.dburi.GetText())
	if err != nil {
		c.addStatus(fmt.Sprintf("Tried opening DB with driver %s but got error: %s", driver, err))
		return
	}

	c.addStatus("Successfully connected to server")
	c.addStatus("Trying to access database...")

	//make sure we actually have access
	if err := db.Ping(); err != nil {
		c.addStatus(fmt.Sprintf("Tried to ping DB but got error %s", err))
		return
	}

	stats := db.Stats()
	c.addStatus("Successfully connected to database")
	c.addStatus(fmt.Sprintf("Current open database connections %d", stats.OpenConnections))
	c.db = db

	//now enable list tables button
}

func (c *Console) listTables() {
	_, driver := c.driver.GetCurrentOption()

	//setup interface
	var dbc dbInter
	if driver == PGSQLOPTION {
		dbc = newPG(c.db)
	}

	if driver == MYSQLOPTION {
		dbc = newMY(c.db)
	}

	dbc = newMY(c.db)
	tables, err := dbc.getTables()
	if err != nil {
		c.addStatus(fmt.Sprintf("Tried to list tables in DB but got error %s", err))
		return
	}

	c.addStatus(fmt.Sprintf("Got %d tables", len(tables)))
	for _, v := range tables {
		c.addStatus(fmt.Sprintf("Table %s", v))
	}
}

func (c *Console) addStatus(logline string) {
	c.logbuf += fmt.Sprintf("%s\n", logline)
	c.status.SetText(c.logbuf)
	c.status.ScrollToEnd()
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

func (c *Console) changeDriver(label string, index int) {
	//check if status is nil - otherwise you get nil pointer
	if c.status == nil {
		return
	}

	//var DB *sql.DB
	c.addStatus(fmt.Sprintf("DB Driver: %s", label))

	if label == PGSQLOPTION {
		//c.dburi.SetText("postgres://<username>:<password>@<host>/<dbname>?sslmode=<verify,disable>")
		//hack for now
		connuri := "postgres://graph:graph@127.0.0.1:5432/graph?sslmode=disable"
		c.dburi.SetText(connuri)
	}

	if label == MYSQLOPTION {
		c.dburi.SetText("<username>:<password>@<host:port>/<dbname>?<paramN=valueN,...>")
	}
	c.app.SetFocus(c.dburi)
}

func (c *Console) DisableButton(label string, disable bool) {
	b := c.form.GetButtonIndex(label)
	c.form.GetButton(b).SetDisabled(disable)
}

func (c *Console) setLayout() {
	c.form.AddDropDown(DB_DRIVER_FIELD, []string{PGSQLOPTION, MYSQLOPTION}, 0, c.changeDriver)
	c.driver = c.form.GetFormItemByLabel(DB_DRIVER_FIELD).(*tview.DropDown)

	c.form.AddInputField(DB_URI_FIELD, "", FIELD_WIDTH, nil, nil)
	c.dburi = c.form.GetFormItemByLabel(DB_URI_FIELD).(*tview.InputField)

	c.form.AddTextArea(DB_QUERY_FIELD, "", FIELD_WIDTH, 10, 500, nil)
	c.query = c.form.GetFormItemByLabel(DB_QUERY_FIELD).(*tview.TextArea)

	c.form.AddTextView(STATUS_FIELD, "", FIELD_WIDTH, 10, true, true)
	c.status = c.form.GetFormItemByLabel(STATUS_FIELD).(*tview.TextView)
	c.status.SetScrollable(true)
	c.status.ScrollToEnd()

	//at startup set default driver
	c.changeDriver(PGSQLOPTION, 0)

	c.form.AddButton(CONNECT_BUTTON, c.Connect)
	c.form.AddButton(TABLES_BUTTON, c.listTables)
	c.form.AddButton(SAVE_BUTTON, c.Save)
	c.form.AddButton(QUIT_BUTTON, c.Close)

	//disable listing tables until we have connection
	c.DisableButton(TABLES_BUTTON, true)
	c.DisableButton(SAVE_BUTTON, true)

}
