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
    X, Y int
    Ativo bool
}

func powerup(jogo *Jogo, comandoCh <-chan AlertaPowerUp, mutex *sync.Mutex, p *PowerUpStruct, guardioesCh []chan AlertaGuardiao) {
    ticker := time.NewTicker(200 * time.Millisecond)
    defer ticker.Stop()

    // Inicialmente coloca o powerup no mapa
    mutex.Lock()
    jogo.Mapa[p.Y][p.X] = PowerUp
    p.Ativo = true
    mutex.Unlock()

    for {
        select {
        case cmd := <-comandoCh:
            mutex.Lock()
            if (cmd.Resgatado || cmd.Destruido) && p.Ativo {
                jogo.Mapa[p.Y][p.X] = Vazio
                p.Ativo = false
                jogo.UltimoVisitado = Vazio 

                if cmd.Resgatado {
                    for _, ch := range guardioesCh {
                        ch <- AlertaGuardiao{Pausado: true}
                    }
                }

                mutex.Unlock()
                return
            }
            mutex.Unlock()

        case <-ticker.C:
            continue
        }
    }
}