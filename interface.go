// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute

// Definições de cores utilizadas no jogo
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

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Renderiza todo o estado atual do jogo na tela
func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	// Desenha todos os elementos do mapa
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			if elem.simbolo == Chave.simbolo && chavePegou {
				interfaceDesenharElemento(x, y, Vazio)
			} else {
				interfaceDesenharElemento(x, y, elem)
			}

		}
	}

	// Desenha o portal se ainda está ativo
	if portalAtivo {
		interfaceDesenharCaractere(26, 22, Portal.simbolo, Portal.cor, Portal.corFundo)
	}

	// Desenha os guardiões por cima
	for _, g := range jogo.Guardiões {
		interfaceDesenharElemento(g.X, g.Y, Inimigo)
	}

	// Desenha o personagem por cima de tudo
	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	// Desenha a barra de status
	interfaceDesenharBarraDeStatus(jogo, chavePegou, portalAtivo, tempoRestante)

	// Atualiza a tela
	interfaceAtualizarTela()
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
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
	interfaceDesenharTextoCentro("Selecione a dificuldade:", 8)
	interfaceDesenharTextoCentro("1 - Facil (40 segundos)", 10)
	interfaceDesenharTextoCentro("2 - Dificil (20 segundos)", 11)
	interfaceDesenharTextoCentro("ESC - Sair", 13)
	interfaceAtualizarTela()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case '1':
				return 40
			case '2':
				return 20
			}
			if ev.Key == termbox.KeyEsc {
				interfaceFinalizar()
				os.Exit(0)
			}
		}
	}
}
