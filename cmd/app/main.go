package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/C0d3-5t3w/GoNetPad/internal/config"
	"github.com/C0d3-5t3w/GoNetPad/internal/locale"
	"github.com/C0d3-5t3w/GoNetPad/internal/logger"
	"github.com/C0d3-5t3w/GoNetPad/internal/ui"
	"github.com/C0d3-5t3w/GoNetPad/internal/ui/themes"
	"github.com/gorilla/websocket"
)

func getHostIP() string {
	var hostIP string
	if envIP := os.Getenv("GONETPAD_HOST"); envIP != "" {
		return envIP
	}

	fmt.Print("Enter host IP address (default: 127.0.0.1): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		hostIP = scanner.Text()
	}
	if hostIP == "" {
		hostIP = "127.0.0.1"
	}
	logger.InfoLogger.Printf("Using host IP: %s\n", hostIP)
	return hostIP
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	hostIP := getHostIP()
	a := app.New()

	if err := locale.Configure(a); err != nil {
		logger.ErrorLogger.Printf("Failed to configure locale: %v\n", err)
	}

	a.Settings().SetTheme(themes.NewBaseTheme())

	w := a.NewWindow(config.AppConfig.WindowTitle)
	w.Resize(fyne.NewSize(float32(config.AppConfig.WindowWidth), float32(config.AppConfig.WindowHeight)))
	w.SetFixedSize(false)

	editor := ui.NewEditor(w)
	var currentText string

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for range ticker.C {
			currentText = editor.TextArea.Text
			broadcast <- currentText
		}
	}()

	http.HandleFunc("/ws", handleWebSocketConnection)
	http.HandleFunc("/", serveIndexHTML)

	go handleBroadcastMessages()

	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New", func() { editor.AddNewTab("Untitled") }),
			fyne.NewMenuItem("Open", func() { openFile(editor, w) }),
			fyne.NewMenuItem("Save", func() { saveFile(editor, w) }),
			fyne.NewMenuItem("Save As...", func() { saveFileAs(editor, w) }),
			fyne.NewMenuItem("Quit", func() { w.Close() }),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Undo", func() { undo(editor) }),
			fyne.NewMenuItem("Redo", func() { redo(editor) }),
			fyne.NewMenuItem("Cut", func() { cut(editor) }),
			fyne.NewMenuItem("Copy", func() { copy(editor) }),
			fyne.NewMenuItem("Paste", func() { paste(editor) }),
			fyne.NewMenuItem("Find", func() { find(editor, w) }),
			fyne.NewMenuItem("Replace", func() { replace(editor, w) }),
		),
		fyne.NewMenu("View",
			fyne.NewMenuItem("Split Horizontally", func() { splitHorizontally(editor) }),
			fyne.NewMenuItem("Split Vertically", func() { splitVertically(editor) }),
			fyne.NewMenuItem("Toggle Line Numbers", func() { toggleLineNumbers(editor) }),
			fyne.NewMenuItem("Toggle Status Bar", func() { toggleStatusBar(editor) }),
		),
		fyne.NewMenu("Themes",
			fyne.NewMenuItem("Base Dark", func() {
				a.Settings().SetTheme(themes.NewBaseTheme())
			}),
		),
	)

	w.SetMainMenu(menu)

	go func() {
		logger.InfoLogger.Printf("Starting WebSocket server on %s%s\n", hostIP, config.AppConfig.WebSocketPort)
		if err := http.ListenAndServe(hostIP+config.AppConfig.WebSocketPort, nil); err != nil {
			logger.ErrorLogger.Printf("WebSocket server error: %v\n", err)
		}
	}()

	go func() {
		logger.InfoLogger.Printf("Starting HTML server on %s%s\n", hostIP, config.AppConfig.HTMLPort)
		if err := http.ListenAndServe(hostIP+config.AppConfig.HTMLPort, nil); err != nil {
			logger.ErrorLogger.Printf("HTML server error: %v\n", err)
		}
	}()

	w.ShowAndRun()
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ErrorLogger.Println("Error upgrading to websocket:", err)
		return
	}
	defer ws.Close()
	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			logger.ErrorLogger.Println("Error reading websocket message:", err)
			delete(clients, ws)
			break
		}
	}
}

func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("pkg/website")).ServeHTTP(w, r)
}

func handleBroadcastMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				logger.ErrorLogger.Println("Error writing websocket message:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func splitHorizontally(editor *ui.Editor) {
	if editor.CurrentView != nil {

		if split, ok := editor.CurrentView.(*container.Split); ok {
			newSplit := container.NewVSplit(split.Leading, split.Trailing)
			editor.CurrentView = newSplit
			editor.Window.SetContent(editor.CurrentView)
		} else {

			entryTextArea := widget.NewMultiLineEntry()
			entryTextArea.SetText(editor.TextArea.Text)
			newEntryTextArea := widget.NewMultiLineEntry()
			newEntryTextArea.SetText(editor.TextArea.Text)
			newSplit := container.NewVSplit(
				container.NewVBox(entryTextArea),
				container.NewVBox(newEntryTextArea),
			)
			editor.CurrentView = newSplit
			editor.Window.SetContent(editor.CurrentView)
		}
	} else {

		logger.InfoLogger.Println("Cannot split: no current view available")
	}
	logger.InfoLogger.Println("Split horizontally action triggered")
}

func splitVertically(editor *ui.Editor) {
	if editor.CurrentView != nil {

		if split, ok := editor.CurrentView.(*container.Split); ok {
			newSplit := container.NewHSplit(split.Leading, split.Trailing)
			editor.CurrentView = newSplit
			editor.Window.SetContent(editor.CurrentView)
		} else {

			entryTextArea := widget.NewMultiLineEntry()
			entryTextArea.SetText(editor.TextArea.Text)
			newEntryTextArea := widget.NewMultiLineEntry()
			newEntryTextArea.SetText(editor.TextArea.Text)
			newSplit := container.NewHSplit(
				container.NewVBox(entryTextArea),
				container.NewVBox(newEntryTextArea),
			)
			editor.CurrentView = newSplit
			editor.Window.SetContent(editor.CurrentView)
		}
	} else {

		logger.InfoLogger.Println("Cannot split: no current view available")
	}
	logger.InfoLogger.Println("Split vertically action triggered")
}

func toggleLineNumbers(editor *ui.Editor) {
	if editor.LineNumbers.Visible {
		editor.LineNumbers.Hide()
	} else {
		editor.LineNumbers.Show()
	}
	logger.InfoLogger.Println("Toggle line numbers action triggered")
}

func toggleStatusBar(editor *ui.Editor) {
	if editor.StatusBar.Visible() {
		editor.StatusBar.Hide()
	} else {
		editor.StatusBar.Show()
	}
	logger.InfoLogger.Println("Toggle status bar action triggered")
}

func cut(editor *ui.Editor) {
	selectedText := editor.GetSelectedText()
	editor.Window.Clipboard().SetContent(selectedText)
	editor.SetText(strings.Replace(editor.GetText(), selectedText, "", 1))
	logger.InfoLogger.Println("Cut action triggered")
}

func copy(editor *ui.Editor) {
	editor.Window.Clipboard().SetContent(editor.GetSelectedText())
	logger.InfoLogger.Println("Copy action triggered")
}

func paste(editor *ui.Editor) {
	content := editor.Window.Clipboard().Content()
	if content != "" {
		if editor.GetSelectedText() != "" {
			editor.SetText(strings.Replace(editor.GetText(), editor.GetSelectedText(), content, 1))
		} else {
			curPos := editor.TextArea.CursorRow*len(editor.GetText()) + editor.TextArea.CursorColumn
			newText := editor.GetText()[:curPos] + content + editor.GetText()[curPos:]
			editor.SetText(newText)
		}
	}
	logger.InfoLogger.Println("Paste action triggered")
}

func find(editor *ui.Editor, w fyne.Window) {
	d := dialog.NewEntryDialog("Find", "Enter text to find:", func(text string) {
		if text == "" {
			return
		}

		index := strings.Index(editor.TextArea.Text, text)
		if index >= 0 {

			editor.StatusBar.ShowTemporaryMessage(fmt.Sprintf("Found at position %d", index))
		} else {
			editor.StatusBar.ShowTemporaryMessage("Text not found")
		}
	}, w)
	d.Show()
	logger.InfoLogger.Println("Find action triggered")
}

func replace(editor *ui.Editor, w fyne.Window) {
	findEntry := widget.NewEntry()
	findEntry.SetPlaceHolder("Find")
	replaceEntry := widget.NewEntry()
	replaceEntry.SetPlaceHolder("Replace with")

	content := container.NewVBox(
		findEntry,
		replaceEntry,
	)

	dialog.ShowCustomConfirm("Replace", "Replace", "Cancel", content, func(confirmed bool) {
		if confirmed && findEntry.Text != "" {
			newText := strings.Replace(editor.TextArea.Text, findEntry.Text, replaceEntry.Text, -1)
			editor.TextArea.SetText(newText)
			count := strings.Count(editor.TextArea.Text, findEntry.Text) - strings.Count(newText, findEntry.Text)
			editor.StatusBar.ShowTemporaryMessage(fmt.Sprintf("Replaced %d occurrences", count))
		}
	}, w)

	logger.InfoLogger.Println("Replace action triggered")
}

func openFile(editor *ui.Editor, w fyne.Window) {
	logger.InfoLogger.Println("Open file action triggered")
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			logger.ErrorLogger.Printf("Error opening file: %v", err)
			dialog.ShowError(err, w)
			return
		}
		if reader == nil {
			logger.InfoLogger.Println("No file selected")
			return
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			logger.ErrorLogger.Printf("Error reading file: %v", err)
			dialog.ShowError(err, w)
			return
		}

		editor.TextArea.SetText(string(data))
		editor.SetFilePath(reader.URI().Path())
		editor.StatusBar.ShowTemporaryMessage(fmt.Sprintf("Opened %s", reader.URI().Name()))
		logger.InfoLogger.Printf("File opened successfully: %s", reader.URI().Path())
	}, w)
	fd.Show()
}

func saveFile(editor *ui.Editor, w fyne.Window) {
	logger.InfoLogger.Println("Save file action triggered")
	if editor.FilePath == "" {
		saveFileAs(editor, w)
		return
	}

	err := os.WriteFile(editor.FilePath, []byte(editor.TextArea.GetText()), 0644)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	editor.StatusBar.ShowTemporaryMessage("File saved")
}

func saveFileAs(editor *ui.Editor, w fyne.Window) {
	logger.InfoLogger.Println("Save as action triggered")
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(editor.TextArea.GetText()))
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		editor.StatusBar.ShowTemporaryMessage("File saved")
		editor.SetFilePath(writer.URI().Path())
	}, w)
	fd.Show()
}

func undo(editor *ui.Editor) {
	logger.InfoLogger.Println("Undo action triggered")
	text, ok := editor.History.Undo()
	if ok {
		editor.TextArea.SetText(text)
	}
}

func redo(editor *ui.Editor) {
	logger.InfoLogger.Println("Redo action triggered")
	text, ok := editor.History.Redo()
	if ok {
		editor.TextArea.SetText(text)
	}
}
