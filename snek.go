package main

import (
	tl "github.com/JoelOtter/termloop"
	"fmt"
	"math/rand"
	"time"
)

const (
	GridSize = 1 // Change to 1 since we won't use it for rendering
	Width    = 40
	Height   = 20
)

type Point struct {
	X, Y int
}

type Snake struct {
	body      []Point
	dx, dy    int
	growing   bool
}

type Food struct {
	pos Point
}

type GameLevel struct {
	*tl.BaseLevel
	snake *Snake
	food  *Food
	score int
}

func NewSnake() *Snake {
	// Create initial body with 10 segments
	initialBody := make([]Point, 10)
	for i := 0; i < 10; i++ {
		initialBody[i] = Point{10 - i, 10}
	}
	
	return &Snake{
		body:   initialBody,
		dx:     1,
		dy:     0,
	}
}

func (s *Snake) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey {
		switch ev.Key {
		case tl.KeyArrowRight:
			if s.dx == 0 { s.dx, s.dy = 1, 0 }
		case tl.KeyArrowLeft:
			if s.dx == 0 { s.dx, s.dy = -1, 0 }
		case tl.KeyArrowUp:
			if s.dy == 0 { s.dx, s.dy = 0, -1 }
		case tl.KeyArrowDown:
			if s.dy == 0 { s.dx, s.dy = 0, 1 }
		}
	}
}

func (s *Snake) Reset() {
	// Create initial body with 10 segments
	initialBody := make([]Point, 10)
	for i := 0; i < 10; i++ {
		initialBody[i] = Point{10 - i, 10}
	}
	s.body = initialBody
	s.dx = 1
	s.dy = 0
	s.growing = false
}

func (s *Snake) CollidesWithSelf() bool {
	head := s.body[0]
	// Start from 1 to skip the head
	for i := 1; i < len(s.body); i++ {
		if head.X == s.body[i].X && head.Y == s.body[i].Y {
			return true
		}
	}
	return false
}

func (s *Snake) Move() {
	head := s.body[0]
	newHead := Point{head.X + s.dx, head.Y + s.dy}

	// Check for wall collision
	if newHead.X < 0 || newHead.Y < 0 || newHead.X >= Width || newHead.Y >= Height {
		s.Reset()
		return
	}

	s.body = append([]Point{newHead}, s.body...)

	if !s.growing {
		s.body = s.body[:len(s.body)-1]
	} else {
		s.growing = false
	}

	// Check for self collision
	if s.CollidesWithSelf() {
		s.Reset()
	}
}

func (s *Snake) Draw(screen *tl.Screen) {
	// Draw body segments (all except head)
	for i := 1; i < len(s.body); i++ {
		p := s.body[i]
		bodyChar := '□'
		if i%2 == 0 {
			bodyChar = '■'
		}
		screen.RenderCell(p.X, p.Y, &tl.Cell{Fg: tl.ColorGreen, Ch: bodyChar})
	}

	// Draw head with direction (in cyan)
	head := s.body[0]
	var headChar rune
	switch {
	case s.dx > 0:
		headChar = '>'
	case s.dx < 0:
		headChar = '<'
	case s.dy > 0:
		headChar = 'v'
	case s.dy < 0:
		headChar = '^'
	}
	screen.RenderCell(head.X, head.Y, &tl.Cell{Fg: tl.ColorCyan, Ch: headChar})
}

func NewFood() *Food {
	rand.Seed(time.Now().UnixNano())
	fx := rand.Intn(Width-2) + 1  // Keep away from borders
	fy := rand.Intn(Height-2) + 1
	return &Food{
		pos: Point{fx, fy},
	}
}

func (f *Food) Draw(screen *tl.Screen) {
	screen.RenderCell(f.pos.X, f.pos.Y, &tl.Cell{Fg: tl.ColorRed, Ch: '★'})
}

func (s *Snake) CollidesWith(p Point) bool {
	head := s.body[0]
	return head.X == p.X && head.Y == p.Y
}

func (f *Food) Respawn() {
	rand.Seed(time.Now().UnixNano())
	f.pos.X = rand.Intn(Width-2) + 1  // Keep away from borders
	f.pos.Y = rand.Intn(Height-2) + 1
}

func (s *Snake) Grow() {
	s.growing = true
}

func (g *GameLevel) Draw(screen *tl.Screen) {
	// Draw snake body segments
	for i := 1; i < len(g.snake.body); i++ {
		p := g.snake.body[i]
		bodyChar := '□'
		if i%2 == 0 {
			bodyChar = '■'
		}
		screen.RenderCell(p.X, p.Y, &tl.Cell{Fg: tl.ColorGreen, Ch: bodyChar})
	}

	// Draw snake head
	head := g.snake.body[0]
	var headChar rune
	switch {
	case g.snake.dx > 0:
		headChar = '>'
	case g.snake.dx < 0:
		headChar = '<'
	case g.snake.dy > 0:
		headChar = 'v'
	case g.snake.dy < 0:
		headChar = '^'
	}
	screen.RenderCell(head.X, head.Y, &tl.Cell{Fg: tl.ColorCyan, Ch: headChar})

	// Draw food
	screen.RenderCell(g.food.pos.X, g.food.pos.Y, &tl.Cell{Fg: tl.ColorRed, Ch: '★'})

	// Draw border
	for x := 0; x < Width; x++ {
		screen.RenderCell(x, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: '─'})
		screen.RenderCell(x, Height-1, &tl.Cell{Fg: tl.ColorWhite, Ch: '─'})
	}
	for y := 0; y < Height; y++ {
		screen.RenderCell(0, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '│'})
		screen.RenderCell(Width-1, y, &tl.Cell{Fg: tl.ColorWhite, Ch: '│'})
	}
	// Draw corners
	screen.RenderCell(0, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: '┌'})
	screen.RenderCell(Width-1, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: '┐'})
	screen.RenderCell(0, Height-1, &tl.Cell{Fg: tl.ColorWhite, Ch: '└'})
	screen.RenderCell(Width-1, Height-1, &tl.Cell{Fg: tl.ColorWhite, Ch: '┘'})

	// Draw score
	scoreText := []rune(fmt.Sprintf("Score: %d", g.score))
	for i, ch := range scoreText {
		screen.RenderCell(2+i, Height, &tl.Cell{Fg: tl.ColorYellow, Ch: ch})
	}
}

func (g *GameLevel) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey {
		switch ev.Key {
		case tl.KeyArrowRight:
			if g.snake.dx == 0 { g.snake.dx, g.snake.dy = 1, 0 }
		case tl.KeyArrowLeft:
			if g.snake.dx == 0 { g.snake.dx, g.snake.dy = -1, 0 }
		case tl.KeyArrowUp:
			if g.snake.dy == 0 { g.snake.dx, g.snake.dy = 0, -1 }
		case tl.KeyArrowDown:
			if g.snake.dy == 0 { g.snake.dx, g.snake.dy = 0, 1 }
		}
	}
}

func main() {
	game := tl.NewGame()
	snake := NewSnake()
	food := NewFood()
	
	level := &GameLevel{
		BaseLevel: tl.NewBaseLevel(tl.Cell{}),
		snake:     snake,
		food:      food,
		score:     0,
	}

	game.Screen().SetLevel(level)
	go func() {
		for {
			head := snake.body[0]
			newHead := Point{head.X + snake.dx, head.Y + snake.dy}

			// Check for wall collision
			if newHead.X <= 0 || newHead.Y <= 0 || newHead.X >= Width-1 || newHead.Y >= Height-1 {
				snake.Reset()
				continue
			}

			snake.body = append([]Point{newHead}, snake.body...)

			if !snake.growing {
				snake.body = snake.body[:len(snake.body)-1]
			} else {
				snake.growing = false
			}

			// Check for self collision
			if snake.CollidesWithSelf() {
				snake.Reset()
				continue
			}

			// Check for food collision
			if newHead.X == food.pos.X && newHead.Y == food.pos.Y {
				snake.growing = true
				food.Respawn()
				level.score++
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
	game.Start()
}