package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/c0d3-5t3w/GoNetPad/internal/capture"
	"github.com/c0d3-5t3w/GoNetPad/internal/config"
	"github.com/c0d3-5t3w/GoNetPad/internal/helpers"
	"github.com/c0d3-5t3w/GoNetPad/internal/locale"
	"github.com/c0d3-5t3w/GoNetPad/internal/logger"
	"github.com/c0d3-5t3w/GoNetPad/internal/themes"
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

var filename string
var updatingContent bool
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

	a.Settings().SetTheme(themes.InvertedTheme)

	w := a.NewWindow("5T3Ws Custom Text Editor")

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for range ticker.C {
			if base64Image, err := capture.CaptureWindow(w); err == nil {
				broadcast <- base64Image
			}
		}
	}()

	textArea := widget.NewMultiLineEntry()
	textArea.SetPlaceHolder("Enter Text Here...")

	go handleWebSocketConnections()
	go handleWebSocketMessages()
	go serveIndexHTML()

	formatButton := widget.NewButton("Format Code", func() {
		logger.InfoLogger.Println("Format Code button clicked")
		formatted, err := helpers.FormatCode(textArea.Text)
		if err != nil {
			dialog.ShowError(err, w)
			logger.ErrorLogger.Println("Error formatting code:", err)
			return
		}
		updatingContent = true
		textArea.SetText(formatted)
		updatingContent = false
		broadcast <- formatted
	})

	lineNumbers := widget.NewLabel("1")

	content := container.NewBorder(nil, container.NewHBox(formatButton), lineNumbers, nil, textArea)
	scroll := container.NewScroll(content)

	updateLineNumbers := func(content string) {
		lines := strings.Split(content, "\n")
		var sb strings.Builder
		for i := 1; i <= len(lines); i++ {
			sb.WriteString(fmt.Sprintf("%d\n", i))
		}
		lineNumbers.SetText(sb.String())
		logger.InfoLogger.Println("Line numbers updated")
	}

	history := binding.NewStringList()
	history.Append(textArea.Text)
	historyIndex := 0

	updateHistory := func() {
		history.Append(textArea.Text)
		historyIndex = history.Length() - 1
		logger.InfoLogger.Println("History updated, index:", historyIndex)
	}

	undo := func() {
		if historyIndex > 0 {
			historyIndex--
			val, _ := history.GetValue(historyIndex)
			updatingContent = true
			textArea.SetText(val)
			updatingContent = false
			updateLineNumbers(val)
			logger.InfoLogger.Println("Undo action performed, index:", historyIndex)
			broadcast <- val
		}
	}

	redo := func() {
		if historyIndex < history.Length()-1 {
			historyIndex++
			val, _ := history.GetValue(historyIndex)
			updatingContent = true
			textArea.SetText(val)
			updatingContent = false
			updateLineNumbers(val)
			logger.InfoLogger.Println("Redo action performed, index:", historyIndex)
			broadcast <- val
		}
	}

	textArea.OnChanged = func(content string) {
		if updatingContent {
			return
		}

		updatingContent = true
		defer func() { updatingContent = false }()

		if last, err := history.GetValue(historyIndex); err == nil && last == content {
			return
		}
		logger.InfoLogger.Println("Text area content changed")
		updateHistory()
		updateLineNumbers(content)
		broadcast <- content
	}

	openFile := func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				logger.ErrorLogger.Println("Error opening file:", err)
				return
			}
			if reader == nil {
				return
			}
			filename = reader.URI().Path()
			file, err := os.Open(filename)
			if err != nil {
				dialog.ShowError(err, w)
				logger.ErrorLogger.Println("Error opening file:", err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			var content strings.Builder
			for scanner.Scan() {
				line := scanner.Text()
				content.WriteString(line + "\n")
			}
			updatingContent = true
			textArea.SetText(content.String())
			updatingContent = false
			updateHistory()
			updateLineNumbers(content.String())
			broadcast <- content.String()
			logger.InfoLogger.Println("File opened:", filename)
		}, w).Show()
	}

	saveFile := func() {
		if filename == "" {
			dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					logger.ErrorLogger.Println("Error saving file:", err)
					return
				}
				if writer == nil {
					return
				}
				filename = writer.URI().Path()
				file, err := os.Create(filename)
				if err != nil {
					dialog.ShowError(err, w)
					logger.ErrorLogger.Println("Error creating file:", err)
					return
				}
				defer file.Close()

				content := textArea.Text
				writer.Write([]byte(content))
				logger.InfoLogger.Println("File saved:", filename)
			}, w).Show()
		} else {
			file, err := os.Create(filename)
			if err != nil {
				dialog.ShowError(err, w)
				logger.ErrorLogger.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			content := textArea.Text
			file.WriteString(content)
			logger.InfoLogger.Println("File saved:", filename)
		}
	}

	newFile := func() {
		filename = ""
		updatingContent = true
		textArea.SetText("")
		updatingContent = false
		updateHistory()
		updateLineNumbers("")
		broadcast <- ""
		logger.InfoLogger.Println("New file created")
	}

	findText := func() {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Enter text to find")

		dialog.ShowCustomConfirm("Find Text", "Find", "Cancel", entry, func(b bool) {
			if !b {
				return
			}
			searchText := entry.Text
			if searchText == "" {
				return
			}
			content := textArea.Text
			index := strings.Index(content, searchText)
			if index == -1 {
				dialog.ShowInformation("Not Found", "Text not found", w)
				return
			}
			textArea.CursorRow = strings.Count(content[:index], "\n")
			textArea.CursorColumn = index - strings.LastIndex(content[:index], "\n") - 1
			textArea.Refresh()
			logger.InfoLogger.Println("Text found:", searchText)
		}, w)
	}

	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("New", func() { newFile() }),
			fyne.NewMenuItem("Open", func() { openFile() }),
			fyne.NewMenuItem("Save", func() { saveFile() }),
		),
		fyne.NewMenu("Edit",
			fyne.NewMenuItem("Undo", func() { undo() }),
			fyne.NewMenuItem("Redo", func() { redo() }),
			fyne.NewMenuItem("Find", func() { findText() }),
		),
		fyne.NewMenu("Themes",
			fyne.NewMenuItem("Inverted Theme", func() {
				a.Settings().SetTheme(themes.InvertedTheme)
			}),
		),
	)

	w.Canvas().AddShortcut(&fyne.ShortcutUndo{}, func(_ fyne.Shortcut) {
		undo()
	})
	w.Canvas().AddShortcut(&fyne.ShortcutRedo{}, func(_ fyne.Shortcut) {
		redo()
	})

	w.SetMainMenu(menu)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(false)

	go func() {
		logger.InfoLogger.Printf("Starting WebSocket server on %s%s\n", hostIP, config.WebSocketPort)
		if err := http.ListenAndServe(hostIP+config.WebSocketPort, nil); err != nil {
			logger.ErrorLogger.Printf("WebSocket server error: %v\n", err)
		}
	}()

	go func() {
		logger.InfoLogger.Printf("Starting HTML server on %s%s\n", hostIP, config.HTMLPort)
		if err := http.ListenAndServe(hostIP+config.HTMLPort, nil); err != nil {
			logger.ErrorLogger.Printf("HTML server error: %v\n", err)
		}
	}()

	w.ShowAndRun()
}

func handleWebSocketConnections() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.ErrorLogger.Println("Error upgrading to websocket:", err)
			return
		}
		defer conn.Close()
		clients[conn] = true

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				logger.ErrorLogger.Println("Error reading websocket message:", err)
				delete(clients, conn)
				break
			}
		}
	})
}

func serveIndexHTML() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("website")).ServeHTTP(w, r)
	})
}

func handleWebSocketMessages() {
	for {
		msg := <-broadcast
		if !strings.HasPrefix(msg, "iVBOR") {
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
}
