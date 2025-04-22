// personagem.go - Funções para movimentação e ações do personagem
package main

import "fmt"

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo, powerUpCh chan<- AlertaPowerUp) {
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1 // Move para cima
	case 'a': dx = -1 // Move para a esquerda
	case 's': dy = 1  // Move para baixo
	case 'd': dx = 1  // Move para a direita
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy

	if jogoPodeMoverPara(jogo, nx, ny) {
		// Se há PowerUp, envia alerta pelo canal do personagem
		if jogo.Mapa[ny][nx].simbolo == PowerUp.simbolo {
			powerUpCh <- AlertaPowerUp{Resgatado: true}
		}
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	// Atualmente apenas exibe uma mensagem de status
	jogo.StatusMsg = fmt.Sprintf("Interagindo em (%d, %d)", jogo.PosX, jogo.PosY)
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo, powerUpCh chan<- AlertaPowerUp) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo, powerUpCh)
	}
	return true // Continua o jogo
}
