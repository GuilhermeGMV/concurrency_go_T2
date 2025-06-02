package main

import (
	"log"
	"net/rpc"
	"time"
)

var meuID int
var jogo Jogo
var clientRPC *rpc.Client
var tempoLimite int
var chavePegou bool
var tempoInicio time.Time

func main() {
	// Conecta no servidor
	var err error
	clientRPC, err = rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("Erro ao conectar no servidor RPC: %v\n", err)
	}
	defer clientRPC.Close()

	// Registra jogador
	err = clientRPC.Call("Servidor.RegistrarJogador", "Cliente", &meuID)
	if err != nil {
		log.Fatalf("Erro em RegistrarJogador: %v\n", err)
	}
	log.Printf("Cliente: meuID = %d\n", meuID)

	// Obtem o mapa estático do servidor
	var mapaReply ReplyGetMapa
	err = clientRPC.Call("Servidor.GetMapa", ArgsGetMapa{}, &mapaReply)
	if err != nil {
		log.Fatalf("Erro em GetMapa: %v\n", err)
	}

	jogo = jogoNovo()
	jogoConfigurarMapaServer(mapaReply.Linhas, &jogo)

	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Iniciar goroutine que escuta teclas do usuário
	eventosCh := make(chan EventoTeclado)
	go func() {
		for {
			ev := interfaceLerEventoTeclado()
			eventosCh <- ev
		}
	}()

	// Loop principal com ticks de 50 ms
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	// Loop principal do jogo
	for {
		select {
		// Quando o usuário pressiona alguma tecla
		case ev := <-eventosCh:
			if continuar := personagemExecutarAcao(ev); !continuar {
				return
			}

		// A cada tick, atualiza o estado e redesenha
		case <-ticker.C:
			// obtém o estado global do servidor
			var estadoReply GetEstadoReply
			err := clientRPC.Call("Servidor.GetEstado", GetEstadoArgs{}, &estadoReply)
			if err != nil {
				log.Printf("Erro em GetEstado: %v\n", err)
				continue
			}

			// Atualiza variáveis globais
			chavePegou = estadoReply.ChavePegou
			if estadoReply.TempoLimite > 0 {
				tempoLimite = estadoReply.TempoLimite
			}
			if estadoReply.ChaveTimestamp > 0 {
				tempoInicio = time.Unix(estadoReply.ChaveTimestamp, 0)
			}

			// Verifica se houve vitória
			if estadoReply.Vitoria {
				finalizarComVitória()
				return
			}

			// Verifica se o jogador ainda está no jogo
			encontrado := false
			for _, j := range estadoReply.Jogadores {
				if j.ID == meuID {
					encontrado = true
					break
				}
			}
			if !encontrado {
				// Jogador morreu
				finalizarComMorte()
				return
			}

			// Atualiza o mapa do servidor
			var mapaReply ReplyGetMapa
			err = clientRPC.Call("Servidor.GetMapa", ArgsGetMapa{}, &mapaReply)
			if err != nil {
				log.Printf("Erro em GetMapa: %v\n", err)
				continue
			}

			jogoConfigurarMapaServer(mapaReply.Linhas, &jogo)

			// desenha mapa estático e, em seguida, todos os jogadores
			interfaceLimparTela()

			// desenha o tabuleiro a partir da linha 1
			for y, linhaElems := range jogo.Mapa {
				for x, elem := range linhaElems {
					interfaceDesenharElemento(x, y+1, elem)
				}
			}

			// desenha guardiões
			for _, g := range estadoReply.Guardioes {
				interfaceDesenharElemento(g.X, g.Y+1, Inimigo)
			}

			// desenha jogadores
			for _, pl := range estadoReply.Jogadores {
				interfaceDesenharElemento(pl.X, pl.Y+1, Personagem)
			}

			// Desenha barra de status com timer global
			if chavePegou && tempoLimite > 0 && tempoInicio.Unix() > 0 {
				tempoRestante := tempoLimite - int(time.Now().Sub(tempoInicio).Seconds())
				if tempoRestante <= 0 {
					finalizarComMorte()
					return
				}
				interfaceDesenharBarraDeStatus(&jogo, chavePegou, true, tempoRestante)
			} else {
				interfaceDesenharBarraDeStatus(&jogo, false, false, 0)
			}

			interfaceAtualizarTela()
		}
	}
}
