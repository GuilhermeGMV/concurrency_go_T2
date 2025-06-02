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

	guard1 := &GuardiaoInterno{X: 20, Y: 10, UltimoVisitado: ' '}
	guard2 := &GuardiaoInterno{X: 30, Y: 12, UltimoVisitado: ' '}
	guard3 := &GuardiaoInterno{X: 20, Y: 20, UltimoVisitado: ' '}
	guard4 := &GuardiaoInterno{X: 40, Y: 26, UltimoVisitado: ' '}
	guard5 := &GuardiaoInterno{X: 11, Y: 21, UltimoVisitado: ' '}
	guard6 := &GuardiaoInterno{X: 65, Y: 27, UltimoVisitado: ' '}

	// Insere no struct do servidor
	server.Guardioes = append(server.Guardioes, guard1, guard2, guard3, guard4, guard5, guard6)

	// canal para notificar guardiões sobre Resgatado/Destruido de power-up
	guardioesPowerUpCh := make(chan AlertaPowerUp)

	// canais individuais para cada guardião (para sinais de Detectado e Pausado)
	ch1 := make(chan AlertaGuardiao)
	ch2 := make(chan AlertaGuardiao)
	ch3 := make(chan AlertaGuardiao)
	ch4 := make(chan AlertaGuardiao)
	ch5 := make(chan AlertaGuardiao)
	ch6 := make(chan AlertaGuardiao)

	// Armazena esses canais em slice para enviar "Pausado" em massa
	guardioesCh := []chan AlertaGuardiao{ch1, ch2, ch3, ch4, ch5, ch6}

	// Inicia goroutine dos guardiões
	go guardiao(server, ch1, &mu, guard1, guardioesPowerUpCh)
	go guardiao(server, ch2, &mu, guard2, guardioesPowerUpCh)
	go guardiao(server, ch3, &mu, guard3, guardioesPowerUpCh)
	go guardiao(server, ch4, &mu, guard4, guardioesPowerUpCh)
	go guardiao(server, ch5, &mu, guard5, guardioesPowerUpCh)
	go guardiao(server, ch6, &mu, guard6, guardioesPowerUpCh)

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

	// Cria canal para sinais de colisão "jogador toca power-up"
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

	// Rotina que verifica se algum jogador entrou na célula da chave
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			mu.Lock()
			for _, pj := range server.Jogo.Jogadores {
				// Verifica se o jogador está na posição da chave
				if []rune(server.Jogo.Mapa[pj.Y])[pj.X] == '🔑' {
					log.Printf("Jogador pegou a chave na posição (%d,%d)\n", pj.X, pj.Y)
					// Remove a chave do mapa
					server.substituirMapa(pj.X, pj.Y, ' ')
					// Cria o portal na posição fixa (14,26)
					server.substituirMapa(14, 26, '⧉')
					log.Printf("Portal criado na posição (14,26)\n")
				}
				// Verifica se o jogador está na posição do portal
				if []rune(server.Jogo.Mapa[pj.Y])[pj.X] == '⧉' {
					log.Printf("Jogador entrou no portal na posição (%d,%d)\n", pj.X, pj.Y)
					// Remove o jogador do mapa (vitória)
					delete(server.Jogo.Jogadores, pj.ID)
					// Marca vitória
					server.Vitoria = true
				}
			}
			mu.Unlock()
		}
	}()

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
