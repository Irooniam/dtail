package internal

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/rivo/tview"
)

//go:embed env_config
var CONFIG string

const (
	DB_DRIVER_FIELD   = "DB driver"
	DB_URI_FIELD      = "DB Connection URI"
	DB_QUERY_FIELD    = "DB query"
	STATUS_FIELD      = "status"
	CONNECT_BUTTON    = "Connect DB"
	TABLES_BUTTON     = "Choose Table"
	TABLES_LIST_FIELD = "Tables"
	SAVE_BUTTON       = "Save"
	QUIT_BUTTON       = "Quit"
	PGSQLOPTION       = "postgres"
	MYSQLOPTION       = "mysql"
	FIELD_WIDTH       = 100
	CONFIG_LISTEN     = "dtail_table_update"
)

type Console struct {
	app    *tview.Application
	form   *tview.Form
	driver *tview.InputField
	dburi  *tview.InputField
	tables *tview.DropDown
	status *tview.TextView
	dbc    dbInter
	logbuf string
}

type dbInter interface {
	getTables() ([]string, error)
	createTriggers(table string) error
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
	driver := c.driver.GetText()
	c.addStatus(fmt.Sprintf("Preparing to connect to database server using %s driver", driver))

	db, err := NewDB(driver, c.dburi.GetText())
	if err != nil {
		c.addStatus(fmt.Sprintf("Tried opening DB with driver %s but got error: %s", driver, err))
		return
	}

	//make sure we actually have access
	if err := db.Ping(); err != nil {
		c.addStatus(fmt.Sprintf("Tried to ping DB but got error %s", err))
		return
	}

	stats := db.Stats()
	c.addStatus("Successfully connected to database")
	c.addStatus(fmt.Sprintf("Current open database connections %d", stats.OpenConnections))

	var dbc dbInter
	if driver == PGSQLOPTION {
		dbc = newPG(db)
	}

	c.dbc = dbc

	//now enable list tables button
	c.listTables()
	c.DisableButton(TABLES_BUTTON, false)
	c.DisableButton(CONNECT_BUTTON, true)
}

func (c *Console) listTables() {
	tables, err := c.dbc.getTables()
	if err != nil {
		c.addStatus(fmt.Sprintf("Tried to list tables in DB but got error %s", err))
		return
	}

	//remove all options and start fresh
	for i := 0; i < c.tables.GetOptionCount(); i++ {
		c.tables.RemoveOption(i)
	}

	c.addStatus(fmt.Sprintf("Found %d tables", len(tables)))
	c.tables.SetDisabled(false)
	for _, v := range tables {
		c.tables.AddOption(v, nil)
	}

	//make sure we have tables before setting default
	if len(tables) > 0 {
		c.tables.SetCurrentOption(0)
	}
}

func (c *Console) saveTable() {
	_, table := c.tables.GetCurrentOption()
	c.addStatus(fmt.Sprintf("Saving configuration for table %s", table))

	out := fmt.Sprintf(CONFIG, c.dburi.GetText(), CONFIG_LISTEN)
	err := os.WriteFile(".env", []byte(out), 0644)
	if err != nil {
		c.addStatus(fmt.Sprintf("Tried saving file to .env but got error: %s", err))
		return
	}

	c.addStatus("Saved configuration file to .env")
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

func (c *Console) chooseTable() {
	_, table := c.tables.GetCurrentOption()
	c.addStatus(fmt.Sprintf("Creating triggers for table %s", table))

	if err := c.dbc.createTriggers(table); err != nil {
		c.addStatus(fmt.Sprintf("Tried to create triggers but got error: %s", err))
	}

	c.addStatus("Triggers created successfully")
	c.DisableButton(TABLES_BUTTON, true)
	c.DisableButton(SAVE_BUTTON, false)
}

func (c *Console) DisableButton(label string, disable bool) {
	b := c.form.GetButtonIndex(label)
	c.form.GetButton(b).SetDisabled(disable)
}

func (c *Console) setLayout() {
	//hardcode driver as we only support pgsql at the mo
	c.form.AddInputField(DB_DRIVER_FIELD, "postgres", FIELD_WIDTH, nil, nil)
	c.driver = c.form.GetFormItemByLabel(DB_DRIVER_FIELD).(*tview.InputField)
	c.driver.SetDisabled(true)

	//connection string
	c.form.AddInputField(DB_URI_FIELD, "", FIELD_WIDTH, nil, nil)
	c.dburi = c.form.GetFormItemByLabel(DB_URI_FIELD).(*tview.InputField)
	c.dburi.SetText("postgres://<username>:<password>@<host>/<dbname>?sslmode=<verify,disable>")

	// ############ hack for now
	c.dburi.SetText("postgres://graph:graph@127.0.0.1/graph?sslmode=disable")

	//tables dropdown
	c.form.AddDropDown(TABLES_LIST_FIELD, []string{"None"}, 0, nil)
	c.tables = c.form.GetFormItemByLabel(TABLES_LIST_FIELD).(*tview.DropDown)
	c.tables.SetDisabled(true)

	c.form.AddTextView(STATUS_FIELD, "", FIELD_WIDTH, 10, true, true)
	c.status = c.form.GetFormItemByLabel(STATUS_FIELD).(*tview.TextView)
	c.status.SetScrollable(true)
	c.status.ScrollToEnd()

	//at startup set default driver

	c.form.AddButton(CONNECT_BUTTON, c.Connect)
	c.form.AddButton(TABLES_BUTTON, c.chooseTable)
	c.form.AddButton(SAVE_BUTTON, c.saveTable)
	c.form.AddButton(QUIT_BUTTON, c.Close)

	//disable listing tables until we have connection
	c.DisableButton(TABLES_BUTTON, true)
	c.DisableButton(SAVE_BUTTON, true)

}
