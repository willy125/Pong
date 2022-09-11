package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF
const PaddleHeight = 4
const initialBallVelocityRow = 1
const initialBallVelocityCol = 2

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var debugLog string

var gameObjects []*GameObject

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

func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()
	for {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()
		time.Sleep(75 * time.Millisecond)

	}
}
func UpdateState() {
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}
}
func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)
	for _, obj := range gameObjects {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()
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
func HandleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < screenHeight {
		player1Paddle.row++
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < screenHeight {
		player2Paddle.row++
	}
}
func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {

			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan

}
func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight, symbol: PaddleSymbol,
		velRow: 0, velCol: 0,
	}
	player2Paddle = &GameObject{
		row: paddleStart, col: width - 1, width: 1, height: PaddleHeight, symbol: PaddleSymbol,
		velRow: 0, velCol: 0,
	}
	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1, symbol: BallSymbol,
		velRow: initialBallVelocityRow, velCol: initialBallVelocityCol,
	}
	gameObjects = []*GameObject{
		player1Paddle, player2Paddle, ball,
	}
}
func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}
