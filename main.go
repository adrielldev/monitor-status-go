package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 2
const delay = 5

func main() {
	exibeIntroducao()
	menu()
}

func menu() {
	for {

		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 3:
			configuraSites()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)

		default:
			fmt.Println("Não conheço este comando")
			os.Exit(-1)

		}
	}
}

func exibeIntroducao() {
	nome := adicionaOuLeConfiguracoes()
	fmt.Println("Olá, sr", nome)
	fmt.Println("Programa na versão 1.2")
}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir logs")
	fmt.Println("3 - Adicionar sites")
	fmt.Println("0 - Sair do programa")

}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("")

	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	sites := leSitesArquivo()

	for i := 0; i < monitoramento; i++ {
		for _, site := range sites {
			testaSite(site)

		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}

}

func testaSite(site string) {

	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, "está com problemas")
		registraLog(site, false)
	}
}

func leSitesArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}
	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString(byte('\n'))
		linha = strings.TrimSpace(linha)
		if err == io.EOF || linha == "" {
			break
		}
		sites = append(sites, linha)

	}

	arquivo.Close()

	return sites

}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04::05") + " - " + site + " - Online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()

}

func imprimeLogs() {
	arquivo, err := os.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)

	}
	fmt.Println(string(arquivo))

}

func adicionaOuLeConfiguracoes() string {
	arquivo, err := os.OpenFile("config.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}
	leitor := bufio.NewReader(arquivo)
	username, _ := leitor.ReadString('\n')
	if username == "" {
		fmt.Println("Você ainda não possui configurações de usuário, digite o nome desejado:")
		var nome string
		fmt.Scan(&nome)
		arquivo.WriteString(nome)
		fmt.Println("Agora para começar a usar a aplicação digite os sites que deseja monitorar:")
		configuraSites()
	}
	arquivo.Close()

	return username
}

func configuraSites() {
	arquivo, err := os.OpenFile("sites.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}
	fmt.Println("No momento você possui os seguintes sites:")
	leSitesArquivo()
	var input string
	fmt.Println("Para encerrar a configuração dos sites digite 0 e dê enter")
	for input != "0" {
		fmt.Scan(&input)
		if input != "0" && input != "" {
			arquivo.WriteString(input + "\n")
		}

	}
	arquivo.Close()
}
