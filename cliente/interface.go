package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type Cor = termbox.Attribute

const (
	CorPadrao      Cor = termbox.ColorDefault
	CorCinzaEscuro     = termbox.ColorDarkGray
	CorVermelho        = termbox.ColorRed
	CorVerde           = termbox.ColorGreen
	CorAmarela         = termbox.ColorYellow
	CorCiano           = termbox.ColorCyan
	CorAzul            = termbox.ColorBlue
	CorParede          = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede     = termbox.ColorDarkGray
	CorTexto           = termbox.ColorDarkGray
)

type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
}

func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyEsc {
				return EventoTeclado{Tipo: "sair"}
			}
			if ev.Ch == 'e' || ev.Ch == 'E' {
				return EventoTeclado{Tipo: "interagir"}
			}
			// se for w/a/s/d → “mover”
			if ev.Ch == 'w' || ev.Ch == 'a' || ev.Ch == 's' || ev.Ch == 'd' ||
				ev.Ch == 'W' || ev.Ch == 'A' || ev.Ch == 'S' || ev.Ch == 'D' {
				return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
			}
			// para outras teclas, continua esperando
		}
	}
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo, chavePegou bool, portalAtivo bool, tempoRestante int) {
	mensagem := "Pegue a chave e vá até o portal!"

	if chavePegou {
		if portalAtivo {
			mensagem = fmt.Sprintf("Vá até o portal! Tempo restante: %ds", tempoRestante)
		} else {
			mensagem = "O portal desapareceu! Você perdeu o tempo!"
		}
	}

	// Limpa linha 0
	largura, _ := interfaceTamanhoTela()
	for x := 0; x < largura; x++ {
		interfaceDesenharCaractere(x, 0, ' ', termbox.ColorDefault, termbox.ColorBlack)
	}

	// Desenha a mensagem
	for i, ch := range mensagem {
		interfaceDesenharCaractere(i, 0, ch, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func interfaceCorFundoVermelho() {
	interfaceCorDeFundo(CorVermelho)
}

func interfaceCorFundoAzul() {
	interfaceCorDeFundo(CorAzul)
}

func interfaceCorFundoAmarelo() {
	interfaceCorDeFundo(CorAmarela)
}

func interfaceCorFundoPreto() {
	interfaceCorDeFundo(CorPadrao)
}

func interfaceCorDeFundo(cor termbox.Attribute) {
	termbox.Clear(CorPadrao, cor)
}

func interfaceDesenharTextoCentro(texto string, linha int) {
	largura, _ := interfaceTamanhoTela()
	colunaInicio := (largura - len(texto)) / 2

	for i, ch := range texto {
		interfaceDesenharCaractere(colunaInicio+i, linha, ch, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func interfaceDesenharCaractere(x, y int, ch rune, corFrente, corFundo termbox.Attribute) {
	termbox.SetCell(x, y, ch, corFrente, corFundo)
}

func interfaceTamanhoTela() (largura, altura int) {
	return termbox.Size()
}

func finalizarComMorte() {
	for range 3 {
		interfaceLimparTela()

		interfaceCorFundoVermelho()
		interfaceAtualizarTela()

		time.Sleep(200 * time.Millisecond)

		interfaceLimparTela()

		interfaceCorFundoPreto()
		interfaceAtualizarTela()

		time.Sleep(200 * time.Millisecond)
	}

	interfaceLimparTela()
	interfaceDesenharTextoCentro("GAME OVER!", 10)
	interfaceAtualizarTela()

	time.Sleep(2 * time.Second)

	interfaceFinalizar()
	os.Exit(0)
}

func finalizarComVitória() {
	for range 3 {
		interfaceLimparTela()

		interfaceCorFundoAzul()
		interfaceAtualizarTela()

		time.Sleep(200 * time.Millisecond)

		interfaceLimparTela()

		interfaceCorFundoAmarelo()
		interfaceAtualizarTela()

		time.Sleep(200 * time.Millisecond)
	}

	interfaceLimparTela()
	interfaceDesenharTextoCentro("Você Venceu!", 10)
	interfaceAtualizarTela()

	time.Sleep(2 * time.Second)

	interfaceFinalizar()
	os.Exit(0)
}

func interfaceSelecionarDificuldade() int {
	interfaceLimparTela()
	writeCentered(10, "1 – Fácil (40 segundos)")
	writeCentered(12, "2 – Difícil (20 segundos)")
	interfaceAtualizarTela()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case '1':
				return 40
			case '2':
				return 20
			case 'E', 'e':
				return 40 // pode tratar ‘e’ como saída de modo simplificado
			}
			if ev.Key == termbox.KeyEsc {
				return 40
			}
		}
	}
}

func writeCentered(y int, msg string) {
	width, _ := termbox.Size()
	x := (width - len(msg)) / 2
	for i, ch := range msg {
		termbox.SetCell(x+i, y, ch, termbox.ColorWhite, termbox.ColorDefault)
	}
}
