package main

import (
	"fmt"
	"os"
	"time"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

const (
	screenWidth  = 40
	screenHeight = 20
	playerChar   = '@'
	zombieChar   = 'Z'
	bulletChar   = '|'
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

type Player struct {
	position Point
	health   int
}

type Zombie struct {
	position Point
	alive    bool
}

type Bullet struct {
	position Point
	active   bool
	dir      int
}

var (
	screen    tcell.Screen
	player    Player
	zombies   []Zombie
	bullets   []Bullet
	gameOver  bool
	score     int
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
	player = Player{
		position: Point{x: screenWidth / 2, y: screenHeight - 2},
		health:   3,
	}
	zombies = []Zombie{}
	for i := 0; i < 5; i++ {
		zombies = append(zombies, Zombie{position: Point{x: rand.Intn(screenWidth - 1), y: rand.Intn(screenHeight / 2)}})
	}
	bullets = []Bullet{}
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

func drawPlayer() {
	screen.SetContent(player.position.x, player.position.y, playerChar, nil, tcell.StyleDefault.Background(tcell.ColorGreen))
}

func drawZombies() {
	for _, z := range zombies {
		if z.alive {
			screen.SetContent(z.position.x, z.position.y, zombieChar, nil, tcell.StyleDefault.Background(tcell.ColorRed))
		}
	}
}

func drawBullets() {
	for _, b := range bullets {
		if b.active {
			screen.SetContent(b.position.x, b.position.y, bulletChar, nil, tcell.StyleDefault.Background(tcell.ColorYellow))
		}
	}
}

func updateZombies() {
	for i := range zombies {
		if zombies[i].alive {
			zombies[i].position.y++
			if zombies[i].position.y >= screenHeight-1 {
				zombies[i].alive = false
				player.health--
			}
		}
	}
}

func updateBullets() {
	for i := range bullets {
		if bullets[i].active {
			switch bullets[i].dir {
			case up:
				bullets[i].position.y--
			}
			if bullets[i].position.y < 1 {
				bullets[i].active = false
			}
		}
	}
}

func checkCollisions() {
	for i := range bullets {
		if bullets[i].active {
			for j := range zombies {
				if zombies[j].alive && bullets[i].position.x == zombies[j].position.x && bullets[i].position.y == zombies[j].position.y {
					bullets[i].active = false
					zombies[j].alive = false
					score++
				}
			}
		}
	}
}

func handleInput() {
	// Handle user input for controlling player and shooting bullets
	switch ev := screen.PollEvent().(type) {
	case *tcell.EventResize:
		screen.Sync()
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			gameOver = true
		case tcell.KeyUp:
			if player.position.y > 1 {
				player.position.y--
			}
		case tcell.KeyDown:
			if player.position.y < screenHeight-2 {
				player.position.y++
			}
		case tcell.KeyLeft:
			if player.position.x > 1 {
				player.position.x--
			}
		case tcell.KeyRight:
			if player.position.x < screenWidth-2 {
				player.position.x++
			}
		case tcell.KeyEnter:
			bullets = append(bullets, Bullet{position: player.position, active: true, dir: up})
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
		drawPlayer()
		drawZombies()
		drawBullets()
		displayScore()

		screen.Show()

		updateZombies()      // Move zombies
		updateBullets()      // Move bullets
		checkCollisions()    // Check for collisions between bullets and zombies
		handleInput()        // Handle player input

		if player.health <= 0 {
			gameOver = true
		}

		time.Sleep(100 * time.Millisecond) // Controls game speed
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
