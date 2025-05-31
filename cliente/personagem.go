// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"log"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(direcao string) {
	var reply MoveReply
	err := clientRPC.Call("Servidor.MoveJogador", MoveArgs{ID: meuID, Direcao: direcao}, &reply)
	if err != nil {
		log.Printf("Erro RPC MoveJogador: %v\n", err)
		return
	}
	// Se reply.Erro não for vazio, significa que o movimento foi bloqueado
	if reply.Erro != "" {
		return
	}
	// Atualiza minha posição local para desenhar o avatar
	jogo.PosX = reply.NovoX
	jogo.PosY = reply.NovoY
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado) bool {
	switch ev.Tipo {
	case "sair":
		return false
	case "interagir":
		// se quiser fazer algo via RPC ao “interagir”, escreva aqui
		return true
	case "mover":
		switch ev.Tecla {
		case 'w', 'W':
			personagemMover("w")
		case 'a', 'A':
			personagemMover("a")
		case 's', 'S':
			personagemMover("s")
		case 'd', 'D':
			personagemMover("d")
		}
	}
	return true
}
