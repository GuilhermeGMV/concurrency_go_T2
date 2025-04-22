// main.go - Loop principal do jogo
package main

import (
	"os"
	"sync"
	"time"
)

var mu sync.Mutex
var chavePegou bool = false
var portalAtivo bool = false
var timerChave *time.Timer
var expiracaoPortal time.Time
var tempoRestante int

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

	powerUpCh := make(chan AlertaPowerUp) // Canal para o personagem
	guardioesPowerUpCh := make(chan AlertaPowerUp) // Canal para todos os guardiões
	powerUp := PowerUpStruct{X: 45, Y: 15}

	// Cria canais individuais para cada guardião
	var canais []chan AlertaGuardiao

	for i := range jogo.Guardiões {
		ch := make(chan AlertaGuardiao)
		canais = append(canais, ch)
		// Guardião envia para o canal dos guardiões
		go guardiao(&jogo, ch, &mu, &jogo.Guardiões[i], guardioesPowerUpCh)
	}

	// PowerUp escuta tanto o canal do personagem quanto dos guardiões
	go powerup(&jogo, powerUpCh, &mu, &powerUp, guardioesPowerUpCh, canais)

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
			if continuar := personagemExecutarAcao(evento, &jogo, powerUpCh); !continuar {
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
				finalizarComMorte()
			}
		}

		// Verifica se jogador pegou a chave
		if !chavePegou && jogo.PosX == 62 && jogo.PosY == 20 {
			chavePegou = true
			portalAtivo = true
			expiracaoPortal = time.Now().Add(40 * time.Second)
			timerChave = time.NewTimer(40 * time.Second)
		}

		// Se a chave foi pega, verificar se o tempo expirou
		if chavePegou && portalAtivo {
			select {
			case <-timerChave.C:
				finalizarComMorte()
			default:
				// Se o timer ainda não acabou, segue normalmente
			}
			tempoRestante = int(time.Until(expiracaoPortal).Seconds())
			if tempoRestante < 0 {
				tempoRestante = 0
			}
		}

		// Verifica se o jogador chegou ao portal a tempo
		if chavePegou && portalAtivo && jogo.PosX == 26 && jogo.PosY == 22 {
			finalizarComVitória()
		}
	}
}
