package main

type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool
}

type Jogo struct {
	Mapa [][]Elemento
	PosX int
	PosY int
}

var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	PowerUp    = Elemento{'★', CorAmarela, CorPadrao, false}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Portal     = Elemento{'⧉', CorCiano, CorPadrao, false}
	Chave      = Elemento{'🔑', CorAmarela, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	return Jogo{
		Mapa: [][]Elemento{},
		PosX: 0,
		PosY: 0,
	}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoConfigurarMapaServer(linhas []string, jogo *Jogo) {
	jogo.Mapa = make([][]Elemento, len(linhas))
	for y, linha := range linhas {
		runes := []rune(linha)
		cols := make([]Elemento, len(runes))
		for x, ch := range runes {
			switch ch {
			case Parede.simbolo:
				cols[x] = Parede
			case Vegetacao.simbolo:
				cols[x] = Vegetacao
			case Chave.simbolo:
				cols[x] = Chave
			case Portal.simbolo:
				cols[x] = Portal
			case PowerUp.simbolo:
				cols[x] = PowerUp
			default:
				cols[x] = Vazio
			}
		}
		jogo.Mapa[y] = cols
	}
}
