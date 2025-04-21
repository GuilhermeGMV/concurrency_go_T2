// main.go - Loop principal do jogo
package main

import (
	"os"
	"sync"
	"time"
)

var mu sync.Mutex

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Cria canais individuais para cada guardião
	var canais []chan AlertaGuardiao
	for i := range jogo.Guardiões {
		ch := make(chan AlertaGuardiao)
		canais = append(canais, ch)
		go guardiao(&jogo, ch, &mu, &jogo.Guardiões[i])
	}

	// Cria canal para eventos de teclado
	eventosCh := make(chan EventoTeclado)

	// Goroutine para ler teclado sem bloquear o programa
	go func() {
		for {
			evento := interfaceLerEventoTeclado()
			eventosCh <- evento
		}
	}()

	// Controle para ativar perseguição apenas uma vez
	ativou := false

	// Loop principal
	for {
		select {
		case evento := <-eventosCh:
			if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
				return
			}

		case <-time.After(50 * time.Millisecond):
			// Só para não travar o loop
		}

		// Lógica de ultrapassagem do x=58
		if !ativou && jogo.PosX > 58 {
			for _, ch := range canais {
				ch <- AlertaGuardiao{Detectado: true}
			}
			ativou = true
		}

		interfaceDesenharJogo(&jogo)

		for _, g := range jogo.Guardiões {
			if g.X == jogo.PosX && g.Y == jogo.PosY {
				finalizarComMorte(&jogo)
			}
		}
	}
}
