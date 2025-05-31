package main

type Jogador struct {
	ID   int
	Nome string
	X    int
	Y    int
}

type ArgsGetMapa struct{}

type ReplyGetMapa struct {
	Linhas []string
}

type MoveArgs struct {
	ID      int    // ID do jogador (recebido em RegistrarJogador)
	Direcao string // “w”, “s”, “a” ou “d”
}

type MoveReply struct {
	NovoX int
	NovoY int
	Erro  string // vazio se tudo OK; string com mensagem se bloqueado
}

type GetEstadoArgs struct{}

type GetEstadoReply struct {
	Jogadores []Jogador
}
