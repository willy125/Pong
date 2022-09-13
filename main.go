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

func PrintStringCenter(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
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

func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()
	for !IsGameOver() {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()
		time.Sleep(75 * time.Millisecond)

	}
	scrrenWidth, screenHeight := screen.Size()
	winner := GetWinner()
	PrintStringCenter(screenHeight/2-1, scrrenWidth/2, "Game Over!")
	PrintStringCenter(screenHeight/2, scrrenWidth/2, fmt.Sprintf("%s 1 wins...", winner))
	screen.Show()
	time.Sleep(3 * time.Second)
	screen.Fini()

}
func UpdateState() {
	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}
	if collidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}
	if CollidesWithPaddle(ball, player1Paddle) || CollidesWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
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
func CollidesWithPaddle(ball *GameObject, paddle *GameObject) bool {
	var collidesColumn bool
	if ball.col < paddle.col {
		collidesColumn = ball.col+ball.velCol >= paddle.col
	} else {
		collidesColumn = ball.col+ball.velCol <= paddle.col
	}
	return collidesColumn &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
}
func collidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}
func IsGameOver() bool {
	return GetWinner() != ""
}
func GetWinner() string {
	screenWidth, _ := screen.Size()
	if ball.col < 0 {
		return "player 2"
	} else if ball.col >= screenWidth {
		return "player1"
	} else {
		return ""
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
