package main

import (
	"math/rand"
	"sync"
	"time"
)

type AlertaGuardiao struct {
	Detectado bool
}

type Guardiao struct {
	X, Y           int
	UltimoVisitado Elemento
}

func guardiao(jogo *Jogo, comandoCh <-chan AlertaGuardiao, mutex *sync.Mutex, g *Guardiao) {
	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	dx, dy := 0, 0
	perseguindo := false

	for {
		select {
		case cmd := <-comandoCh:
			if cmd.Detectado {
				perseguindo = true
			}
		case <-ticker.C:
			if perseguindo {
				mutex.Lock()
				if jogo.PosX > g.X {
					dx, dy = 1, 0
				} else if jogo.PosX < g.X {
					dx, dy = -1, 0
				} else if jogo.PosY > g.Y {
					dx, dy = 0, 1
				} else if jogo.PosY < g.Y {
					dx, dy = 0, -1
				}
				mutex.Unlock()
			} else {
				dx, dy = direcaoAleatoria()
			}
			guardiaoMover(jogo, g, dx, dy, mutex)
		}
	}
}
func guardiaoMover(jogo *Jogo, g *Guardiao, dx, dy int, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()

	nx, ny := g.X+dx, g.Y+dy

	// 1. Primeiro tenta na direção ideal
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogo.Mapa[g.Y][g.X] = g.UltimoVisitado
		g.UltimoVisitado = jogo.Mapa[ny][nx]
		jogo.Mapa[ny][nx] = Inimigo
		g.X, g.Y = nx, ny
		return
	}

	// 2. Se não conseguiu, tenta pequenos desvios aleatórios
	for tentativas := 0; tentativas < 4; tentativas++ {
		adx, ady := direcaoAleatoria()
		nx, ny = g.X+adx, g.Y+ady

		if jogoPodeMoverPara(jogo, nx, ny) {
			jogo.Mapa[g.Y][g.X] = g.UltimoVisitado
			g.UltimoVisitado = jogo.Mapa[ny][nx]
			jogo.Mapa[ny][nx] = Inimigo
			g.X, g.Y = nx, ny
			return
		}
	}

	// 3. Se nada deu certo, fica parado
}

func direcaoAleatoria() (dx, dy int) {
	switch rand.Intn(4) {
	case 0:
		return 1, 0 // Direita
	case 1:
		return -1, 0 // Esquerda
	case 2:
		return 0, 1 // Baixo
	case 3:
		return 0, -1 // Cima
	}
	return 0, 0
}

func distancia(aX, aY, bX, bY int) int {
	dx := aX - bX
	if dx < 0 {
		dx = -dx
	}
	dy := aY - bY
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
