package main

import (
	"sync"
	"time"
)

type AlertaPowerUp struct {
	Resgatado bool
	Destruido bool
}

type PowerUpStruct struct {
	X, Y  int
	Ativo bool
}

func powerup(
	jogo *Servidor,
	personagemCh <-chan AlertaPowerUp,
	mutex *sync.Mutex,
	p *PowerUpStruct,
	guardioesPowerUpCh <-chan AlertaPowerUp,
	guardioesCh []chan AlertaGuardiao,
) {
	// Insere o power-up no mapa do servidor
	mutex.Lock()
	jogo.substituirMapa(p.X, p.Y, '★')
	p.Ativo = true
	mutex.Unlock()

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-personagemCh:
			if (cmd.Resgatado || cmd.Destruido) && p.Ativo {
				mutex.Lock()
				jogo.substituirMapa(p.X, p.Y, ' ')
				p.Ativo = false
				mutex.Unlock()

				if cmd.Resgatado {
					// Avisa aos guardiões para pausarem
					for _, ch := range guardioesCh {
						ch <- AlertaGuardiao{Pausado: true}
					}
				}
				return
			}
		case cmd := <-guardioesPowerUpCh:
			if (cmd.Resgatado || cmd.Destruido) && p.Ativo {
				mutex.Lock()
				jogo.substituirMapa(p.X, p.Y, ' ')
				p.Ativo = false
				mutex.Unlock()
				return
			}

		case <-ticker.C:
			// continua até receber sinal de personagemCh ou guardioesPowerUpCh
		}
	}
}
