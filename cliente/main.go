package main

import (
	"log"
	"net/rpc"
	"time"
)

var meuID int
var jogo Jogo
var clientRPC *rpc.Client

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
	interfaceSelecionarDificuldade()
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

	// Loop principal
	for {
		select {
		// Quando o usuário pressiona alguma tecla
		case ev := <-eventosCh:
			if continuar := personagemExecutarAcao(ev); !continuar {
				return
			}

		// A cada tick, chama GetEstado + redesenha mapa + jogadores
		case <-ticker.C:
			// obtém a posição de todos os jogadores conectados
			var estadoReply GetEstadoReply
			err := clientRPC.Call("Servidor.GetEstado", GetEstadoArgs{}, &estadoReply)
			if err != nil {
				log.Printf("Erro em GetEstado: %v\n", err)
				continue
			}

			// desenha mapa estático e, em seguida, todos os jogadores
			interfaceLimparTela()

			// desenha o tabuleiro
			for y, linhaElems := range jogo.Mapa {
				for x, elem := range linhaElems {
					if elem.simbolo == Chave.simbolo {
						interfaceDesenharElemento(x, y, Vazio)
					} else {
						interfaceDesenharElemento(x, y, elem)
					}
				}
				_ = y
			}

			// desenha guardiões
			for _, g := range estadoReply.Guardioes {
				interfaceDesenharElemento(g.X, g.Y, Inimigo)
			}

			// desenha jogadores
			for _, pl := range estadoReply.Jogadores {
				interfaceDesenharElemento(pl.X, pl.Y, Personagem)
			}

			interfaceAtualizarTela()
		}
	}
}
