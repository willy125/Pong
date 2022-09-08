package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const PaddleSymbol = 0x2588
const PaddleHeight = 4

var screen tcell.Screen
var player1 *Paddle
var player2 *Paddle
var debugLog string

type Paddle struct {
	row, col, width, height int
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		col += 1
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
	}
}
func Print(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)
	Print(player1.row, player1.col, player1.width, player1.height, PaddleSymbol)
	Print(player2.row, player2.col, player2.width, player2.height, PaddleSymbol)
	screen.Show()
}

func main() {
	InitScreen()
	InitGameState()
	InitUserInput()
	inputChan := InitUserInput()
	for {
		DrawState()
		time.Sleep(50 * time.Millisecond)

		key := <-inputChan
		if key == "Rune[q]" {
			screen.Fini()
			os.Exit()
		}
	}
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

}
func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {

			//time.Sleep(75 * time.Millisecond)
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				debugLog = ev.Name()
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan

}
func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1 = &Paddle{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight,
	}
	player2 = &Paddle{
		row: paddleStart, col: width - 1, width: 1, height: PaddleHeight,
	}

}
