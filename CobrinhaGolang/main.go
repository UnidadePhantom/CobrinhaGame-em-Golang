package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// Essas constantes definem o tamanho do tabuleiro e da janela.
	larguraGrade = 20
	alturaGrade  = 20
	tamanhoBloco = 24

	// Essa constante controla a velocidade da cobrinha.
	framesPorPasso = 8

	// Essas constantes ajudam a posicionar os textos na tela.
	larguraTela = larguraGrade * tamanhoBloco
	alturaTela  = alturaGrade*tamanhoBloco + 60
)

var (
	// Essas cores deixam o desenho mais legível e fácil de alterar depois.
	corFundo     = color.RGBA{18, 18, 18, 255}
	corGrade     = color.RGBA{35, 35, 35, 255}
	corCabeca    = color.RGBA{80, 220, 120, 255}
	corCorpo     = color.RGBA{40, 160, 80, 255}
	corComida    = color.RGBA{220, 70, 70, 255}
	corRodape    = color.RGBA{25, 25, 25, 255}
	corSeparador = color.RGBA{55, 55, 55, 255}
)

// Ponto representa uma posição da grade.
type Ponto struct {
	X int
	Y int
}

// Jogo guarda todo o estado necessário para a partida.
type Jogo struct {
	cobra      []Ponto
	direcao    Ponto
	proximaDir Ponto
	comida     Ponto
	frameAtual int
	pontos     int
	gameOver   bool
}

// novoJogo cria uma nova instância já pronta para jogar.
func novoJogo() *Jogo {
	jogo := &Jogo{}
	jogo.reiniciar()
	return jogo
}

// reiniciar coloca a cobrinha no centro e cria uma nova comida.
func (j *Jogo) reiniciar() {
	centroX := larguraGrade / 2
	centroY := alturaGrade / 2

	j.cobra = []Ponto{
		{X: centroX, Y: centroY},
		{X: centroX - 1, Y: centroY},
		{X: centroX - 2, Y: centroY},
	}
	j.direcao = Ponto{X: 1, Y: 0}
	j.proximaDir = j.direcao
	j.frameAtual = 0
	j.pontos = 0
	j.gameOver = false
	j.gerarComida()
}

// Update roda várias vezes por segundo e cuida da lógica do jogo.
func (j *Jogo) Update() error {
	// Se o jogo acabou, apertar R reinicia a partida.
	if j.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			j.reiniciar()
		}
		return nil
	}

	// Essa parte lê o teclado e guarda a próxima direção desejada.
	j.lerTeclado()

	// Esse contador faz a cobrinha andar em passos, e não a cada frame.
	j.frameAtual++
	if j.frameAtual < framesPorPasso {
		return nil
	}
	j.frameAtual = 0

	// Aqui a nova direção passa a valer de verdade.
	j.direcao = j.proximaDir

	// Calcula a nova posição da cabeça da cobra.
	cabecaAtual := j.cobra[0]
	novaCabeca := Ponto{
		X: cabecaAtual.X + j.direcao.X,
		Y: cabecaAtual.Y + j.direcao.Y,
	}

	// Encostar na parede encerra a partida.
	if novaCabeca.X < 0 || novaCabeca.X >= larguraGrade || novaCabeca.Y < 0 || novaCabeca.Y >= alturaGrade {
		j.gameOver = true
		return nil
	}

	// Encostar no próprio corpo também encerra a partida.
	if j.cobraContem(novaCabeca) {
		j.gameOver = true
		return nil
	}

	// A nova cabeça entra no começo da lista.
	novaCobra := append([]Ponto{novaCabeca}, j.cobra...)

	// Se comer a comida, cresce e soma pontos.
	if novaCabeca == j.comida {
		j.cobra = novaCobra
		j.pontos++
		j.gerarComida()
		return nil
	}

	// Se não comer, a cauda sai para manter o mesmo tamanho.
	j.cobra = novaCobra[:len(novaCobra)-1]
	return nil
}

// lerTeclado troca a direção sem permitir giro instantâneo de 180 graus.
func (j *Jogo) lerTeclado() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if j.direcao.Y != 1 {
			j.proximaDir = Ponto{X: 0, Y: -1}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		if j.direcao.Y != -1 {
			j.proximaDir = Ponto{X: 0, Y: 1}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if j.direcao.X != 1 {
			j.proximaDir = Ponto{X: -1, Y: 0}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if j.direcao.X != -1 {
			j.proximaDir = Ponto{X: 1, Y: 0}
		}
	}
}

// gerarComida sorteia uma posição que não esteja ocupada pela cobra.
func (j *Jogo) gerarComida() {
	for {
		novaComida := Ponto{
			X: rand.Intn(larguraGrade),
			Y: rand.Intn(alturaGrade),
		}

		if !j.cobraContem(novaComida) {
			j.comida = novaComida
			return
		}
	}
}

// cobraContem verifica se uma posição já está ocupada pela cobra.
func (j *Jogo) cobraContem(p Ponto) bool {
	for _, parte := range j.cobra {
		if parte == p {
			return true
		}
	}
	return false
}

// Draw desenha tudo o que aparece na janela.
func (j *Jogo) Draw(tela *ebiten.Image) {
	// Primeiro, pinta o fundo geral.
	tela.Fill(corFundo)

	// Depois desenha a grade do tabuleiro.
	for y := 0; y < alturaGrade; y++ {
		for x := 0; x < larguraGrade; x++ {
			ebitenutil.DrawRect(
				tela,
				float64(x*tamanhoBloco),
				float64(y*tamanhoBloco),
				float64(tamanhoBloco-1),
				float64(tamanhoBloco-1),
				corGrade,
			)
		}
	}

	// Desenha a comida.
	ebitenutil.DrawRect(
		tela,
		float64(j.comida.X*tamanhoBloco),
		float64(j.comida.Y*tamanhoBloco),
		float64(tamanhoBloco-1),
		float64(tamanhoBloco-1),
		corComida,
	)

	// Desenha a cobra, com a cabeça em outra cor para destacar.
	for i, parte := range j.cobra {
		corAtual := corCorpo
		if i == 0 {
			corAtual = corCabeca
		}

		ebitenutil.DrawRect(
			tela,
			float64(parte.X*tamanhoBloco),
			float64(parte.Y*tamanhoBloco),
			float64(tamanhoBloco-1),
			float64(tamanhoBloco-1),
			corAtual,
		)
	}

	// Essa faixa inferior funciona como painel de informações.
	ebitenutil.DrawRect(
		tela,
		0,
		float64(alturaGrade*tamanhoBloco),
		float64(larguraTela),
		60,
		corRodape,
	)
	ebitenutil.DrawRect(
		tela,
		0,
		float64(alturaGrade*tamanhoBloco),
		float64(larguraTela),
		2,
		corSeparador,
	)

	// Mostra instruções simples para facilitar o estudo do código.
	mensagem := fmt.Sprintf("Pontos: %d | Mova com WASD ou setas", j.pontos)
	ebitenutil.DebugPrintAt(tela, mensagem, 12, alturaGrade*tamanhoBloco+12)
	ebitenutil.DebugPrintAt(tela, "Se perder, aperte R para reiniciar", 12, alturaGrade*tamanhoBloco+32)

	// Mostra a mensagem de fim de jogo por cima do tabuleiro.
	if j.gameOver {
		ebitenutil.DrawRect(tela, float64(larguraTela/2-110), float64(alturaTela/2-28), 220, 50, color.RGBA{0, 0, 0, 120})
		ebitenutil.DebugPrintAt(tela, "Fim de jogo!", larguraTela/2-45, alturaTela/2-20)
		ebitenutil.DebugPrintAt(tela, "Aperte R para jogar de novo", larguraTela/2-95, alturaTela/2)
	}
}

// Layout informa para a biblioteca o tamanho lógico da janela.
func (j *Jogo) Layout(_, _ int) (int, int) {
	return larguraTela, alturaTela
}

func main() {
	// Semente aleatória para a comida aparecer em lugares diferentes.
	rand.Seed(time.Now().UnixNano())

	// Cria o jogo e configura a janela.
	ebiten.SetWindowSize(larguraTela, alturaTela)
	ebiten.SetWindowTitle("Cobrinha em Go")

	// RunGame inicia o loop principal do jogo.
	if err := ebiten.RunGame(novoJogo()); err != nil {
		log.Fatal(err)
	}
}

// aviso importante ainda estou aprendendo essa linguágem!
