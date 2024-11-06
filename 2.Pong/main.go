package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
)

const (
	paddleHeight = 4
	paddleWidth  = 1
	ballWidth    = 1
	winScore     = 3
)

var (
	screen               tcell.Screen
	ballX, ballY         int
	ballDX, ballDY       int = 1, 1
	leftPaddleY          int
	rightPaddleY         int
	scoreLeft, scoreRight int
	gameOver             bool
	winner               string
)

// Function to print strings on screen
func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

// Display score at the top of the screen
func displayScore(s tcell.Screen) {
	scoreStr := fmt.Sprintf("%d - %d", scoreLeft, scoreRight)
	w, _ := s.Size()
	emitStr(s, w/2-len(scoreStr)/2, 1, tcell.StyleDefault.Foreground(tcell.ColorYellow), scoreStr)
}

// Initialize the game: place paddles and ball in starting positions
func initGame() {
	w, h := screen.Size()
	ballX, ballY = w/2, h/2
	leftPaddleY, rightPaddleY = h/2-paddleHeight/2, h/2-paddleHeight/2
	scoreLeft, scoreRight = 0, 0
	gameOver = false
}

// Display "Game Over" screen with the winner's message
func displayGameOver() {
	screen.Clear()
	w, h := screen.Size()
	gameOverStr := fmt.Sprintf("%s Wins!", winner)
	emitStr(screen, w/2-len(gameOverStr)/2, h/2, tcell.StyleDefault.Foreground(tcell.ColorGreen), gameOverStr)
	emitStr(screen, w/2-9, h/2+2, tcell.StyleDefault, "Press ESC to exit or R to restart.")
	screen.Show()
}

// Draw a paddle at a given x, y position
func drawPaddle(x, y int) {
	for i := 0; i < paddleHeight; i++ {
		screen.SetContent(x, y+i, ' ', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
	}
}

// Draw the ball at its current position
func drawBall() {
	screen.SetContent(ballX, ballY, ' ', nil, tcell.StyleDefault.Background(tcell.ColorRed))
}

// Update the ball's position and handle collision with paddles or walls
func updateBallPosition() {
	if gameOver {
		return
	}

	w, h := screen.Size()
	ballX += ballDX
	ballY += ballDY

	// Ball collision with top and bottom walls
	if ballY <= 0 || ballY >= h-1 {
		ballDY *= -1
	}

	// Ball collision with left paddle
	if ballX == 1 && ballY >= leftPaddleY && ballY < leftPaddleY+paddleHeight {
		ballDX *= -1
	}

	// Ball collision with right paddle
	if ballX == w-2 && ballY >= rightPaddleY && ballY < rightPaddleY+paddleHeight {
		ballDX *= -1
	}

	// Ball goes out of bounds (scoring)
	if ballX <= 0 {
		scoreRight++
		ballX, ballY = w/2, h/2
		ballDX = 1
		checkGameOver()
	} else if ballX >= w-1 {
		scoreLeft++
		ballX, ballY = w/2, h/2
		ballDX = -1
		checkGameOver()
	}
}

// Check if either player has reached the winning score
func checkGameOver() {
	if scoreLeft >= winScore {
		gameOver = true
		winner = "Player 1"
		displayGameOver()
	} else if scoreRight >= winScore {
		gameOver = true
		winner = "Player 2"
		displayGameOver()
	}
}

// Handle player input for controlling paddles and game actions
func handleInput() {
	_, h := screen.Size()
	switch ev := screen.PollEvent().(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			screen.Fini()
			os.Exit(0)
		case tcell.KeyUp:
			if rightPaddleY > 0 {
				rightPaddleY--
			}
		case tcell.KeyDown:
			if rightPaddleY+paddleHeight < h {
				rightPaddleY++
			}
		case tcell.KeyRune:
			if ev.Rune() == 'w' && leftPaddleY > 0 {
				leftPaddleY--
			} else if ev.Rune() == 's' && leftPaddleY+paddleHeight < h {
				leftPaddleY++
			} else if ev.Rune() == 'r' && gameOver {
				initGame()
			}
		}
	}
}

// Main function initializes and runs the game loop
func main() {
	encoding.Register()

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

	initGame()

	// Game loop
	go func() {
		for {
			if !gameOver {
				screen.Clear()
				w, _ := screen.Size()
				drawPaddle(0, leftPaddleY)
				drawPaddle(w-1, rightPaddleY)
				drawBall()
				displayScore(screen)
				screen.Show()

				updateBallPosition()
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// Input handling loop
	for {
		handleInput()
	}
}
