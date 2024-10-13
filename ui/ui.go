package ui

import (
	"fmt"
	"lazytask/tasks"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	app       *tview.Application
	pages     *tview.Pages
	taskTable *tview.Table
	footer    *tview.TextView
}

func InitializeApp() *App {
	appInstance := &App{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
	}
	return appInstance
}

func (a *App) Run() error {
	a.buildUI()
	return a.app.Run()
}

func (a *App) buildUI() {
	a.taskTable = tview.NewTable().
		SetSelectable(true, false).
		SetBorders(false)

	a.footer = tview.NewTextView().
		SetText("Controls: [green]a[white]-Add [green]d[white]-Delete [green]Enter[white]-Edit [green]q[white]-Quit").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	a.loadTasks()

	a.taskTable.SetSelectedFunc(func(row int, col int) {
		task := tasks.GetTasks()[row-1]
		a.showTaskForm(task, false)
	})

	a.taskTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			a.showTaskForm(&tasks.Task{}, true)
		case 'd':
			row, _ := a.taskTable.GetSelection()
			task := tasks.GetTasks()[row-1]
			tasks.DeleteTask(task.ID)
			a.loadTasks()
		case 'q':
			a.app.Stop()
		}
		return event
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.taskTable, 0, 1, true).
		AddItem(a.footer, 1, 0, false)

	a.pages.AddPage("main", layout, true, true)
	a.app.SetRoot(a.pages, true)
}

func (a *App) loadTasks() {
	a.taskTable.Clear()
	headers := [3]string{"ID", "Title", "Status"}

	for i, header := range headers {
		a.taskTable.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	tasksList := tasks.GetTasks()

	for i, task := range tasksList {
		a.taskTable.SetCell(i+1, 0,
			tview.NewTableCell(strconv.Itoa(task.ID)).
				SetTextColor(tcell.ColorWhite))

		a.taskTable.SetCell(i+1, 1,
			tview.NewTableCell(task.Title).
				SetTextColor(tcell.ColorWhite))

		a.taskTable.SetCell(i+1, 2,
			tview.NewTableCell(tasks.StatusCodeToString(task.StatusCode)).
				SetTextColor(tcell.ColorWhite))
	}
}

func (a *App) showTaskForm(task *tasks.Task, isNew bool) {

	statusOptions := []string{"TODO", "IN PROGRESS", "DONE"}

	titleInput := tview.NewInputField().
		SetLabel("Title").
		SetFieldWidth(40).
		SetText(task.Title)

	descriptionInput := tview.NewInputField().
		SetLabel("Description").
		SetFieldWidth(40).
		SetText(task.Description)

	statusDropdown := tview.NewDropDown().
		SetLabel("Status").
		SetOptions(statusOptions, nil).
		SetCurrentOption(task.StatusCode)

	form := tview.NewForm().
		AddFormItem(titleInput).
		AddFormItem(descriptionInput).
		AddFormItem(statusDropdown).
		AddButton("Save", func() {
			task.Title = titleInput.GetText()
			task.Description = descriptionInput.GetText()
			statusIndex, _ := statusDropdown.GetCurrentOption()
			task.StatusCode = statusIndex

			if isNew {
				tasks.AddTask(task)
			} else {
				tasks.UpdateTask(task)
			}
			a.loadTasks()
			a.pages.RemovePage("modal")
			a.app.SetFocus(a.taskTable)
		}).
		AddButton("Cancel", func() {
			a.pages.RemovePage("modal")
			a.app.SetFocus(a.taskTable)
		})

	form.SetBorder(true).
		SetTitle(func() string {
			if isNew {
				return "Add Task"
			}
			return fmt.Sprintf("Edit Task ID: %d", task.ID)
		}()).
		SetTitleAlign(tview.AlignCenter)

	form.SetFocus(0)

	a.pages.AddPage("modal", a.createModal(form, 60, 18), true, true)
	a.app.SetFocus(form)

}

func (a *App) createModal(p tview.Primitive, width, height int) tview.Primitive {
	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false),
			width, 1, true).
		AddItem(nil, 0, 1, false)
	return modal
}
