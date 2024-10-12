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

	a.loadTasks()

	a.taskTable.SetSelectedFunc(func(row int, col int) {
		task := tasks.GetTasks()[row-1]
		a.showTaskModal(task)
	})

	a.taskTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			a.showAddTaskModal()
		case 'd':
			row, _ := a.taskTable.GetSelection()
			task := tasks.GetTasks()[row-1]
			tasks.DeleteTask(task.ID)
			a.loadTasks()
		}
		return event
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.taskTable, 0, 1, true)

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

func (a *App) showTaskModal(task *tasks.Task) {
	statusOptions := []string{"TODO", "IN PROGRESS", "DONE"}
	form := tview.NewForm().
		AddInputField("Title", task.Title, 20, nil, func(text string) {
			task.Title = text
		}).
		AddInputField("Description", task.Description, 20, nil, func(text string) {
			task.Description = text
		}).
		AddDropDown("Status", statusOptions, task.StatusCode, func(option string, index int) {
			task.StatusCode = index
		}).
		AddButton("Save", func() {
			tasks.UpdateTask(task)
			a.loadTasks()
			a.pages.RemovePage("modal")
		}).
		AddButton("Cancel", func() {
			a.pages.RemovePage("modal")
		})

	form.SetBorder(true).SetTitle(fmt.Sprintf("Task ID: %d", task.ID)).SetTitleAlign(tview.AlignCenter)

	a.pages.AddPage("modal", form, true, true)
}

func (a *App) showAddTaskModal() {
	newTask := &tasks.Task{
		Title:       "",
		Description: "",
		StatusCode:  0,
	}

	form := tview.NewForm().
		AddInputField("Title", "", 20, nil, func(text string) {
			newTask.Title = text
		}).
		AddInputField("Description", "", 40, nil, func(text string) {
			newTask.Description = text
		}).
		AddButton("Save", func() {
			tasks.AddTask(newTask)
			a.loadTasks()
			a.pages.RemovePage("modal")
		}).
		AddButton("Cancel", func() {
			a.pages.RemovePage("modal")
		})

	form.SetBorder(true).SetTitle("Add Task").SetTitleAlign(tview.AlignLeft)

	a.pages.AddPage("modal", form, true, true)
	a.app.SetFocus(form)
}
