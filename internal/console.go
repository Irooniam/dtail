package internal

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Console struct {
	App     *tview.Application
	Grid    *tview.Grid
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
	grid := tview.NewGrid()
	app.SetRoot(grid, true).EnableMouse(true).EnablePaste(true)

	return &Console{
		App:  app,
		Grid: grid,
	}

}

func (c *Console) SetLayout() {
	c.Grid.SetBorder(true)
	c.Grid = c.Grid.SetRows(-10, -5, -5, -5, -5, -5, -7, -10, -48)
	c.Grid = c.Grid.SetColumns(-5, -15, -80)

	//header
	header := tview.NewTextView().SetText("Header").SetTextAlign(1)
	c.Grid.AddItem(header, 0, 0, 1, 3, 0, 0, false)

	//host
	host := tview.NewTextView().SetText("host")
	ihost := tview.NewInputField()
	c.Grid.AddItem(host, 1, 0, 1, 1, 0, 0, false)
	c.Grid.AddItem(ihost, 1, 1, 1, 1, 0, 0, false)
	c.Grid.AddItem(tview.NewTextView(), 1, 2, 1, 1, 0, 0, false)

	//user
	user := tview.NewTextView().SetText("user")
	iuser := tview.NewInputField()
	c.Grid.AddItem(user, 2, 0, 1, 1, 0, 0, false)
	c.Grid.AddItem(iuser, 2, 1, 1, 1, 0, 0, false)
	c.Grid.AddItem(tview.NewTextView(), 2, 2, 1, 1, 0, 0, false)

	//password
	passwd := tview.NewTextView().SetText("passwd")
	ipasswd := tview.NewInputField()
	c.Grid.AddItem(passwd, 3, 0, 1, 1, 0, 0, false)
	c.Grid.AddItem(ipasswd, 3, 1, 1, 1, 0, 0, false)
	c.Grid.AddItem(tview.NewTextView(), 3, 2, 1, 1, 0, 0, false)

	//dbname
	dbname := tview.NewTextView().SetText("dbname")
	idbname := tview.NewInputField()
	c.Grid.AddItem(dbname, 4, 0, 1, 1, 0, 0, false)
	c.Grid.AddItem(idbname, 4, 1, 1, 1, 0, 0, false)
	c.Grid.AddItem(tview.NewTextView(), 4, 2, 1, 1, 0, 0, false)

	//tablename
	tables := tview.NewTextView().SetText("tables")
	itables := tview.NewDropDown().SetOptions([]string{"Must first connect to DB    "}, nil).SetCurrentOption(0)
	c.Grid.AddItem(tables, 5, 0, 1, 1, 0, 0, false)
	c.Grid.AddItem(itables, 5, 1, 1, 1, 0, 0, false)
	c.Grid.AddItem(tview.NewTextView(), 5, 2, 1, 1, 0, 0, false)

	//connect button
	button := tview.NewButton("Connect")
	button.SetBorder(false).SetBorderColor(tcell.Color20).SetTitleColor(tcell.ColorRed).SetBorderPadding(1, 1, 1, 1)
	c.Grid.AddItem(button, 6, 1, 1, 1, 0, 0, false)
	c.Grid.AddItem(tview.NewTextView(), 6, 0, 1, 1, 0, 0, false)
	c.Grid.AddItem(tview.NewTextView(), 6, 2, 1, 1, 0, 0, false)

	//spacer
	c.Grid.AddItem(tview.NewTextView(), 7, 0, 1, 3, 0, 0, false)

	//footer
	footer := tview.NewTextView().SetText("footer").SetTextAlign(1)
	c.Grid.AddItem(footer, 8, 0, 1, 3, 0, 0, false)

}
