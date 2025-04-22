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

func powerup(jogo *Jogo, comandoCh <-chan AlertaPowerUp, mutex *sync.Mutex, p *PowerUpStruct) {
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
                // Garante que o PowerUp seja removido do mapa
                jogo.Mapa[p.Y][p.X] = Vazio
                p.Ativo = false
                jogo.UltimoVisitado = Vazio // Importante: atualiza o último visitado
                mutex.Unlock()
                return
            }
            mutex.Unlock()

        case <-ticker.C:
            mutex.Lock()
            if p.Ativo {
                atual := jogo.Mapa[p.Y][p.X]
                if atual == PowerUp {
                    jogo.Mapa[p.Y][p.X] = Vazio
                } else if atual == Vazio {
                    jogo.Mapa[p.Y][p.X] = PowerUp
                }
            }
            mutex.Unlock()
        }
    }
}