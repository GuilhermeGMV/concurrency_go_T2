package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

func main() {
	// Carrega o mapa do arquivo "mapa.txt" e cria o Servidor
	server, err := NewServidor("mapa.txt")
	if err != nil {
		log.Fatalf("Falha ao iniciar servidor: %v\n", err)
	}

	// Registra o Servidor para ser chamado via RPC
	if err := rpc.Register(server); err != nil {
		log.Fatalf("Erro ao registrar RPC: %v\n", err)
	}

	// Inicia os guardiões e o power-up em goroutines separadas:
	var mu sync.Mutex

	guard1 := Guardiao{X: 20, Y: 10, UltimoVisitado: ' '}
	guard2 := Guardiao{X: 30, Y: 12, UltimoVisitado: ' '}

	// canal para notificar guardiões sobre Resgatado/Destruido de power-up
	guardioesPowerUpCh := make(chan AlertaPowerUp)

	// canais individuais para cada guardião (para sinais de Detectado e Pausado)
	ch1 := make(chan AlertaGuardiao)
	ch2 := make(chan AlertaGuardiao)

	// Armazena esses canais em slice para enviar “Pausado” em massa
	guardioesCh := []chan AlertaGuardiao{ch1, ch2}

	// Inicia goroutine dos guardiões
	go guardiao(server, ch1, &mu, &guard1, guardioesPowerUpCh)
	go guardiao(server, ch2, &mu, &guard2, guardioesPowerUpCh)

	// Exemplo de rotina que detecta quando um jogador ultrapassa X > 58 e notifica guardiões
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			mu.Lock()
			for _, pj := range server.Jogo.Jogadores {
				if pj.X > 58 {
					// notifica todos guardiões para entrarem em perseguição
					for _, ch := range guardioesCh {
						ch <- AlertaGuardiao{Detectado: true}
					}
					break
				}
			}
			mu.Unlock()
		}
	}()

	// Cria canal para sinais de colisão “jogador toca power-up”
	personagemToPowerCh := make(chan AlertaPowerUp)
	// power-up em (45, 15)
	power := PowerUpStruct{X: 45, Y: 15, Ativo: false}
	// Rotina que verifica se algum jogador entrou na célula do power-up
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			mu.Lock()
			if power.Ativo {
				for _, pj := range server.Jogo.Jogadores {
					if pj.X == power.X && pj.Y == power.Y {
						personagemToPowerCh <- AlertaPowerUp{Resgatado: true}
						break
					}
				}
			}
			mu.Unlock()
		}
	}()
	// Inicia goroutine do power-up
	go powerup(server, personagemToPowerCh, &mu, &power, guardioesPowerUpCh, guardioesCh)

	// 4) Inicia listener RPC
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Falha ao criar listener: %v\n", err)
	}
	defer listen.Close()
	log.Println("Servidor RPC ouvindo em :1234")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v\n", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
