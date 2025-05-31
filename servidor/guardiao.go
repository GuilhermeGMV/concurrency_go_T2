package main

import (
	"math/rand"
	"sync"
	"time"
)

type AlertaGuardiao struct {
	Detectado bool
	Pausado   bool
}

type Guardiao struct {
	X, Y           int
	UltimoVisitado rune
}

func guardiao(jogo *Servidor, comandoCh <-chan AlertaGuardiao, mutex *sync.Mutex, g *Guardiao, guardioesPowerUpCh chan<- AlertaPowerUp) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	perseguindo := false
	pausado := false

	for {
		select {
		case cmd := <-comandoCh:
			if cmd.Detectado {
				perseguindo = true
			}
			if cmd.Pausado {
				pausado = true
			}
		case <-ticker.C:
			if pausado {
				continue
			}
			dx, dy := 0, 0
			if perseguindo {
				mutex.Lock()
				var alvoX, alvoY int
				for _, pj := range jogo.Jogo.Jogadores {
					alvoX, alvoY = pj.X, pj.Y
					break
				}
				mutex.Unlock()

				if alvoX > g.X {
					dx = 1
				} else if alvoX < g.X {
					dx = -1
				} else if alvoY > g.Y {
					dy = 1
				} else if alvoY < g.Y {
					dy = -1
				}
			} else {
				dx, dy = direcaoAleatoria()
			}

			guardiaoMover(jogo, g, dx, dy, mutex, guardioesPowerUpCh)
		}
	}
}

func guardiaoMover(jogo *Servidor, g *Guardiao, dx, dy int, mutex *sync.Mutex, guardioesPowerUpCh chan<- AlertaPowerUp) {
	mutex.Lock()
	defer mutex.Unlock()

	newX := g.X + dx
	newY := g.Y + dy

	// Verifica limites verticais
	if newY < 0 || newY >= len(jogo.Jogo.Mapa) {
		tentaDesvio(jogo, g, mutex, guardioesPowerUpCh)
		return
	}
	// Verifica limites horizontais
	lineRunes := []rune(jogo.Jogo.Mapa[newY])
	if newX < 0 || newX >= len(lineRunes) {
		tentaDesvio(jogo, g, mutex, guardioesPowerUpCh)
		return
	}

	targetRune := lineRunes[newX]

	// Se encontrar power-up (★), notifica e remove do mapa
	if targetRune == '★' {
		guardioesPowerUpCh <- AlertaPowerUp{Destruido: true}
		jogo.substituirMapa(newX, newY, ' ')
		return
	}

	// Se for parede ('▤') ou outro guardião ('☠'), não se move
	if targetRune == '▤' || targetRune == '☠' {
		return
	}

	// Movimento válido: restaura o que estava em (g.X, g.Y) e move o guardião
	jogo.substituirMapa(g.X, g.Y, g.UltimoVisitado)
	// Guarda o rune que estava na posição destino antes de ocupá-la
	runesNew := []rune(jogo.Jogo.Mapa[newY])
	g.UltimoVisitado = runesNew[newX]
	// Coloca o símbolo do guardião no mapa
	jogo.substituirMapa(newX, newY, '☠')
	g.X = newX
	g.Y = newY

	// Verifica colisão com jogador (se algum jogador ocupa a mesma posição)
	for id, pj := range jogo.Jogo.Jogadores {
		if pj.X == newX && pj.Y == newY {
			// Marca jogador como “morto” removendo do mapa de jogadores
			delete(jogo.Jogo.Jogadores, id)
		}
	}
}

func direcaoAleatoria() (int, int) {
	r := rand.Intn(4)
	switch r {
	case 0:
		return 1, 0 // direita
	case 1:
		return -1, 0 // esquerda
	case 2:
		return 0, 1 // baixo
	default:
		return 0, -1 // cima
	}
}

func (s *Servidor) substituirMapa(x, y int, novoRune rune) {
	s.mu.Lock()
	defer s.mu.Unlock()

	linha := []rune(s.Jogo.Mapa[y])
	linha[x] = novoRune
	s.Jogo.Mapa[y] = string(linha)
}

func tentaDesvio(
	server *Servidor,
	g *Guardiao,
	mutex *sync.Mutex,
	guardioesPowerUpCh chan<- AlertaPowerUp,
) {
	for i := 0; i < 4; i++ {
		adx, ady := direcaoAleatoria()
		tx := g.X + adx
		ty := g.Y + ady

		// Verifica novamente limites
		if ty < 0 || ty >= len(server.Jogo.Mapa) {
			continue
		}
		runesTy := []rune(server.Jogo.Mapa[ty])
		if tx < 0 || tx >= len(runesTy) {
			continue
		}

		tr := runesTy[tx]
		if tr == '▤' || tr == '☠' {
			continue
		}
		// Se for power-up
		if tr == '★' {
			guardioesPowerUpCh <- AlertaPowerUp{Destruido: true}
			server.substituirMapa(tx, ty, ' ')
			return
		}
		// Movimento alternativo válido
		server.substituirMapa(g.X, g.Y, g.UltimoVisitado)
		g.UltimoVisitado = tr
		server.substituirMapa(tx, ty, '☠')
		g.X = tx
		g.Y = ty
		// Verifica colisão com jogador
		for id, pj := range server.Jogo.Jogadores {
			if pj.X == tx && pj.Y == ty {
				delete(server.Jogo.Jogadores, id)
			}
		}
		return
	}
}
