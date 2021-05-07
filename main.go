package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"
)

type Block struct {
	data  []byte
	nonce int /* Precisa existir um dado dentro do seu bloco que vai mudar,
	pra permitir procurar um hash que atenda os requisitos. */
	hash []byte
}

func main() {
	var myBlock Block

	// Definindo a dificuldade de se encontrar um hash (MINERAR)
	const dificuldade = 24
	target := big.NewInt(1)
	target.Lsh(target, uint(256-dificuldade))
	//

	// Formatando os numeros pra ficar bonito de ver na tela, o requisito pra mineracao
	// Esse codigo nao faz diferenca pro entendimento de blockchain
	// Eh soh formatacao na tela
	fmt.Printf("Requisito da mineracao (target) baseado na dificuldade %d\n", dificuldade)
	fmt.Println("** Seu hash precisa ser menor que isso **")
	fmt.Println("Em decimal.")
	fmt.Println(target)
	fmt.Println()
	fmt.Println("Mesmo numero, em binario")
	var b string
	for i := 1; i < dificuldade; i++ {
		b = b + "0"
	}
	b = b + fmt.Sprintf("%b", target)
	fmt.Println(b)
	fmt.Println()
	// ====

	// Pegando o texto, dado, transacao - qquer merda - que vamos colocar no bloco
	fmt.Println("Digita o texto que vc quer incluir na blockchain:")
	reader := bufio.NewReader(os.Stdin)
	data, _ := reader.ReadBytes('\n')
	myBlock.data = data
	fmt.Println("Dados para inserir no blockchain: " + string(myBlock.data))
	// Os dados da transacao estao em myBlock.data
	//

	// Para inserir no blockchain, eu preciso de um hash (como se fosse um endereco)
	// Vamos gerar o hash combinando o myBlock.data com um numero.
	// No caso do bitcoin, mistura com mais coisa. Deixei mais simples, mas mantendo o conceito.
	// EH ISSO QUE EH A MINERACAO. ENCONTRAR UM HASH QUE ATENDA O REQUISITO
	fmt.Println("Procurando hash (MINERANDO)...")
	start := time.Now()

	// Eu comeco com o nonce em zero
	nonce := 0
	var intHash big.Int
	for nonce < math.MaxInt64 {
		// Junto o dado de myBlock.data com myBlock.nonce
		data = bytes.Join([][]byte{myBlock.data, toHex(int64(nonce))}, []byte{})
		hash := sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(target) == -1 {
			// Se entrou aqui, eh pq o hash gerado atendeu os requisitos de dificuldade
			// Bora salvar tudo pra checar depois.
			myBlock.nonce = nonce
			myBlock.hash = hash[:]
			break
		} else {
			// Pra cada tentativa que a mistura de myBlock.data + myBlock.nonce nao atende a dificuldade
			// incrementa nonce
			nonce++
		}
	}

	elapsed := time.Now().Sub(start)

	// Se eu cheguei aqui, terminei a mineracao
	// Encontrei um hash que atende o requisito de dificuldade

	// Formatando os numeros pra ficar bonito de ver na tela, o requisito pra mineracao
	// Esse codigo nao faz diferenca pro entendimento de blockchain
	// Eh soh formatacao na tela

	fmt.Printf("MINERACAO CONCLUIDA\nHash encontrado apos %d tentativas em %d milissegundos\n", myBlock.nonce, elapsed.Milliseconds())
	fmt.Println()
	fmt.Println("O seguinte hash aparentemente atendeu os requisitos de dificuldade da mineracao")
	fmt.Println("** Lembra. Tem de ser menor que o TARGET")

	fmt.Println(intHash.SetBytes(myBlock.hash[:]))

	fmt.Println()
	fmt.Println("Mesmo numero, em binario")
	bh := fmt.Sprintf("%b", intHash.SetBytes(myBlock.hash[:]))
	b = ""
	for i := 0; i < 256-len(bh); i++ {
		b = b + "0"
	}

	b = b + bh

	fmt.Println(b)
	fmt.Println()
	// ====

	// Uma vez o hash encontrado, isso precisa ser checado. Imagina. Qquer um poderia
	// inventar qquer merda.
	//
	// O grande barato do hash:
	//     - eh dificil de se encontrar um que atenda a dificuldade
	//     - mas eh facil de checar se o hash eh valido.
	fmt.Println("Checando o hash")
	start = time.Now()

	data = bytes.Join([][]byte{myBlock.data, toHex(int64(myBlock.nonce))}, []byte{})
	hash := sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	if intHash.Cmp(target) == -1 {
		// Se entrou aqui, eh pq o hash gerado eh valido
		elapsed = time.Now().Sub(start)
		fmt.Printf("Hash VALIDO. Checagem efetuada em %d FUCKING *MICRO* segundo(s)", elapsed.Microseconds())

	} else {
		// Se chegou aqui, tem algum problema com o hash
		log.Panicln("O HASH NAO EH VALIDO")
	}
}

// ToHex from int64 to []byte
func toHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
