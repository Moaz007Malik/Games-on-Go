package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

const (
	screenWidth  = 40
	screenHeight = 20
	snakeChar    = 'O'
	foodChar     = 'X'
)

// Directions
const (
	up = iota
	down
	left
	right
)

type Point struct {
	x, y int
}

type Snake struct {
	body []Point
	dir  int
	grow bool
}

var (
	screen   tcell.Screen
	snake    Snake
	food     Point
	gameOver bool
	score    int
)

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

func displayScore() {
	scoreStr := fmt.Sprintf("Score: %d", score)
	w, _ := screen.Size()
	emitStr(screen, w/2-len(scoreStr)/2, 1, tcell.StyleDefault.Foreground(tcell.ColorYellow), scoreStr)
}

func initGame() {
	// Snake starts with a length of 3
	snake = Snake{
		body: []Point{{x: 10, y: 10}, {x: 9, y: 10}, {x: 8, y: 10}},
		dir:  right,
	}
	food = Point{x: randPos(), y: randPos()}
	gameOver = false
	score = 0
}

func drawBorders() {
	// Draw the borders on all sides of the screen
	for y := 0; y < screenHeight; y++ {
		screen.SetContent(0, y, '|', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
		screen.SetContent(screenWidth-1, y, '|', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
	}
	for x := 0; x < screenWidth; x++ {
		screen.SetContent(x, 0, '-', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
		screen.SetContent(x, screenHeight-1, '-', nil, tcell.StyleDefault.Background(tcell.ColorWhite))
	}
}

func drawSnake() {
	// Draw the snake's body
	for _, p := range snake.body {
		screen.SetContent(p.x, p.y, snakeChar, nil, tcell.StyleDefault.Background(tcell.ColorGreen))
	}
}

func drawFood() {
	// Draw the food item
	screen.SetContent(food.x, food.y, foodChar, nil, tcell.StyleDefault.Background(tcell.ColorRed))
}

func updateSnakePosition() {
	head := snake.body[0]
	var newHead Point

	// Move the snake based on the current direction
	switch snake.dir {
	case up:
		newHead = Point{x: head.x, y: head.y - 1}
	case down:
		newHead = Point{x: head.x, y: head.y + 1}
	case left:
		newHead = Point{x: head.x - 1, y: head.y}
	case right:
		newHead = Point{x: head.x + 1, y: head.y}
	}

	// Insert the new head at the front of the body
	snake.body = append([]Point{newHead}, snake.body...)

	// If the snake is not growing, remove the tail
	if !snake.grow {
		snake.body = snake.body[:len(snake.body)-1]
	} else {
		snake.grow = false
	}

	// Check for collisions with walls or itself
	if newHead.x <= 0 || newHead.x >= screenWidth-1 || newHead.y <= 0 || newHead.y >= screenHeight-1 {
		gameOver = true
	}
	for _, p := range snake.body[1:] {
		if p.x == newHead.x && p.y == newHead.y {
			gameOver = true
		}
	}

	// Check if snake eats food
	if newHead.x == food.x && newHead.y == food.y {
		snake.grow = true
		score++
		// Generate new food position
		food = Point{x: randPos(), y: randPos()}
	}
}

func randPos() int {
	// Ensures that food will never be placed at the borders
	return 1 + rand.Intn(screenWidth-2) // Ensures the food stays within bounds
}

func handleInput() {
	// Handle user input for controlling snake direction
	switch ev := screen.PollEvent().(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			gameOver = true
		case tcell.KeyUp:
			if snake.dir != down {
				snake.dir = up
			}
		case tcell.KeyDown:
			if snake.dir != up {
				snake.dir = down
			}
		case tcell.KeyLeft:
			if snake.dir != right {
				snake.dir = left
			}
		case tcell.KeyRight:
			if snake.dir != left {
				snake.dir = right
			}
		}
	}
}

func displayGameOver() {
	w, h := screen.Size()
	gameOverStr := "Game Over!"
	emitStr(screen, w/2-len(gameOverStr)/2, h/2, tcell.StyleDefault.Foreground(tcell.ColorRed), gameOverStr)
	emitStr(screen, w/2-9, h/2+2, tcell.StyleDefault, "Press ESC to exit.")
	emitStr(screen, w/2-9, h/2+3, tcell.StyleDefault, fmt.Sprintf("Final Score: %d", score))
	screen.Show()
}

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
	for !gameOver {
		screen.Clear()

		drawBorders()
		drawSnake()
		drawFood()
		displayScore()

		screen.Show()

		updateSnakePosition() // Automatically move the snake
		handleInput()          // Allow the user to change direction

		time.Sleep(100 * time.Millisecond) // Controls snake speed
	}

	// Show game over screen
	displayGameOver()

	// Wait for the user to press escape to exit
	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}
