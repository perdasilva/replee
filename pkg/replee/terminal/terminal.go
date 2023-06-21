package terminal

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type Output struct {
	IsErr       bool
	IsSyntaxErr bool
	Output      string
}

type RepleeTerminal struct {
	*tview.Flex
	inputField          *tview.TextArea
	app                 *tview.Application
	commandHistory      []string
	lineHistory         []lineHistoryEntry
	lineIndex           int
	commandHistoryIndex int
	curCommand          string
	totalHeight         int
	execute             func(string) *Output
	onChange            func()
	currentIndent       int
	scrollOffset        int
}

const (
	lineTypePrompt = "prompt"
	lineTypeInput  = "input"
	lineTypeOutput = "output"
)

type lineHistoryEntry struct {
	lineType string
	text     string
}

const prompt = "[yellow]replee:> "
const notPrompt = "         "
const maxLines = 100

func NewRepleeTerminal(app *tview.Application, repl func(string) *Output) *RepleeTerminal {
	inputField := tview.NewTextArea().SetWrap(false).SetLabel(prompt)

	out := &RepleeTerminal{
		Flex:                tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true),
		app:                 app,
		inputField:          inputField,
		commandHistoryIndex: 0,
		lineIndex:           0,
		execute:             repl,
		onChange:            func() {},
	}
	inputField.SetInputCapture(out.handleKeyPush)
	inputField.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		return action, nil
	})
	out.Flex.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseScrollUp {
			out.lineIndex -= 1
			out.scrollOffset = -1
		} else if action == tview.MouseScrollDown {
			out.lineIndex += 1
			out.scrollOffset = 1
		}
		return action, event
	})
	out.AddItem(inputField, 0, 1, true)
	return out
}

func (r *RepleeTerminal) AddItem(item tview.Primitive, fixedSize, proportion int, focus bool) *RepleeTerminal {
	r.Flex.AddItem(item, fixedSize, proportion, focus)
	r.totalHeight += 1
	return r
}

func (r *RepleeTerminal) RemoveItem(item tview.Primitive) *RepleeTerminal {
	r.Flex.RemoveItem(item)
	r.totalHeight -= 1
	return r
}

func (r *RepleeTerminal) addLineHistory(line lineHistoryEntry) {
	r.lineHistory = append(r.lineHistory, line)
	if len(r.lineHistory) > maxLines {
		r.lineHistory = r.lineHistory[1:]
	}
}

func (r *RepleeTerminal) Draw(screen tcell.Screen) {
	r.render()
	r.Flex.Draw(screen)
}

func (r *RepleeTerminal) render() {
	_, _, _, maxRows := r.Flex.GetRect()
	numInputLines := len(strings.Split(r.inputField.GetText(), "\n"))
	maxRows = maxRows - numInputLines
	if r.lineIndex < 0 {
		r.lineIndex = 0
		return
	}
	if r.lineIndex+maxRows > len(r.lineHistory) {
		r.lineIndex = len(r.lineHistory) - maxRows
		if r.lineIndex < 0 {
			r.lineIndex = 0
		}
	}
	r.Clear()
	if maxRows < len(r.lineHistory) && r.scrollOffset == 0 {
		r.lineIndex = len(r.lineHistory) - maxRows
	}
	for i := r.lineIndex; i <= r.lineIndex+maxRows && i <= len(r.lineHistory); i++ {
		if i == len(r.lineHistory) {
			r.AddItem(r.inputField, 0, 1, true)
		} else {
			line := r.lineHistory[i]
			commandView := tview.NewTextView().
				SetChangedFunc(r.onChange).
				SetDynamicColors(true).
				SetText(line.text)

			if line.lineType == lineTypePrompt {
				commandView.SetLabel(prompt)
			}
			r.AddItem(commandView, 1, 1, false)
		}
	}
	r.scrollOffset = 0
}

func (r *RepleeTerminal) handleKeyPush(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyRune && event.Rune() == ' ' && event.Modifiers() == tcell.ModAlt {
		return nil
	}

	if event.Key() == tcell.KeyUp {
		o, _, _, _ := r.inputField.GetCursor()
		if o > 0 {
			return event
		}

		if r.commandHistoryIndex == 0 {
			r.curCommand = r.inputField.GetText()
		}
		if r.commandHistoryIndex+1 <= len(r.commandHistory) {
			r.inputField.SetText(r.commandHistory[len(r.commandHistory)-r.commandHistoryIndex-1], true)
			r.commandHistoryIndex += 1
		}
		return nil
	}
	if event.Key() == tcell.KeyDown {
		o, _, _, _ := r.inputField.GetCursor()
		numLines := len(strings.Split(r.inputField.GetText(), "\n"))

		if numLines != (o+1) && numLines > 1 {
			return event
		}

		if r.commandHistoryIndex == 0 {
			return nil
		}
		if r.commandHistoryIndex-1 > 0 {
			r.inputField.SetText(r.commandHistory[len(r.commandHistory)-r.commandHistoryIndex+1], true)
		} else if r.commandHistoryIndex-1 == 0 {
			r.inputField.SetText(r.curCommand, true)
		}
		r.commandHistoryIndex -= 1
	}
	if event.Key() == tcell.KeyEnter {
		command := r.inputField.GetText()
		if command == "clear" {
			r.lineHistory = []lineHistoryEntry{}
			return nil
		}
		if strings.TrimSpace(command) == "" {
			return nil
		}
		result := r.execute(command)
		if result.IsSyntaxErr {
			indenters := []string{"{", "(", "[", "."}
			for _, indenter := range indenters {
				if strings.HasSuffix(command, indenter) {
					r.currentIndent += 2
					indent := ""
					for i := 0; i < r.currentIndent; i++ {
						indent += " "
					}
					r.inputField.SetText(fmt.Sprintf("%s\n%s", command, indent), true)
					return nil
				}
			}
			r.currentIndent -= 2
			if r.currentIndent < 0 {
				r.currentIndent = 0
			}
			r.inputField.SetText(fmt.Sprintf("%s\n", command), true)
			return nil
		}
		color := "[violet]"
		if result.IsErr {
			color = "[red]"
		}
		output := strings.Trim(result.Output, " \n\t")
		if r.commandHistoryIndex == 0 {
			r.commandHistory = append(r.commandHistory, command)
		}
		r.commandHistoryIndex = 0
		r.currentIndent = 0

		r.RemoveItem(r.inputField)
		r.inputField.SetText("", true)
		for index, line := range strings.Split(command, "\n") {
			if index == 0 {
				r.addLineHistory(lineHistoryEntry{
					lineType: lineTypePrompt,
					text:     line,
				})
			} else {
				r.addLineHistory(lineHistoryEntry{
					lineType: lineTypeInput,
					text:     line,
				})
			}
		}
		for _, line := range strings.Split(output, "\n") {
			if line != "" {
				r.addLineHistory(lineHistoryEntry{
					lineType: lineTypeOutput,
					text:     fmt.Sprintf("%s%s[white]", color, line),
				})
			}
		}
		return nil
	}
	return event
}
