package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"
)

type Jogador struct {
	ID   int
	Nome string
	X    int
	Y    int
}

type EstadoJogo struct {
	Mapa      []string        // cada linha do mapa é uma string
	Jogadores map[int]Jogador // mapeia jogadorID → Jogador{ID,Nome,X,Y}
}

type Servidor struct {
	mu     sync.Mutex
	Jogo   EstadoJogo
	ProxID int
}

type ArgsGetMapa struct{}

type ReplyGetMapa struct {
	Linhas []string
}

type MoveArgs struct {
	ID      int    // ID do jogador
	Direcao string // "w", "s", "a", "d"
}

type MoveReply struct {
	NovoX int
	NovoY int
	Erro  string
}

type GetEstadoArgs struct{}

type GetEstadoReply struct {
	Jogadores []Jogador
}

func NewServidor(mapaPath string) (*Servidor, error) {
	server := &Servidor{
		Jogo: EstadoJogo{
			Mapa:      []string{},
			Jogadores: make(map[int]Jogador),
		},
		ProxID: 1,
	}
	if err := server.carregarMapa(mapaPath); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Servidor) carregarMapa(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var linhas []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		linhas = append(linhas, line)
		if err == io.EOF {
			break
		}
	}

	s.mu.Lock()
	s.Jogo.Mapa = linhas
	s.mu.Unlock()
	return nil
}

func (s *Servidor) RegistrarJogador(nome string, reply *int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.ProxID
	s.ProxID++

	// Posição inicial fixa: (5,5)
	novo := Jogador{ID: id, Nome: nome, X: 5, Y: 5}
	s.Jogo.Jogadores[id] = novo
	*reply = id
	return nil
}

func (s *Servidor) GetMapa(args ArgsGetMapa, reply *ReplyGetMapa) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := make([]string, len(s.Jogo.Mapa))
	copy(c, s.Jogo.Mapa)
	reply.Linhas = c
	return nil
}

func (s *Servidor) MoveJogador(args MoveArgs, reply *MoveReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	j, existe := s.Jogo.Jogadores[args.ID]
	if !existe {
		reply.Erro = "Jogador não encontrado"
		return errors.New("ID inválido")
	}

	dx, dy := 0, 0
	switch args.Direcao {
	case "w":
		dy = -1
	case "s":
		dy = 1
	case "a":
		dx = -1
	case "d":
		dx = 1
	default:
		reply.Erro = "Direção inválida"
		return errors.New("direção inválida")
	}

	newX := j.X + dx
	newY := j.Y + dy

	// Verifica limites verticais
	if newY < 0 || newY >= len(s.Jogo.Mapa) {
		reply.Erro = "Fora dos limites (vertical)"
		return nil
	}
	// Verifica limites horizontais
	if newX < 0 || newX >= len([]rune(s.Jogo.Mapa[newY])) {
		reply.Erro = "Fora dos limites (horizontal)"
		return nil
	}

	// Verifica colisão com parede
	tile := []rune(s.Jogo.Mapa[newY])[newX]
	if tile == '▤' {
		reply.Erro = "Movimento bloqueado: parede"
		return nil
	}

	j.X = newX
	j.Y = newY
	s.Jogo.Jogadores[args.ID] = j

	reply.NovoX = newX
	reply.NovoY = newY
	reply.Erro = ""
	return nil
}

func (s *Servidor) GetEstado(args GetEstadoArgs, reply *GetEstadoReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lista := make([]Jogador, 0, len(s.Jogo.Jogadores))
	for _, j := range s.Jogo.Jogadores {
		lista = append(lista, j)
	}
	reply.Jogadores = lista
	return nil
}
