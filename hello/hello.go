package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 3
const delay = 5
const ARQUIVO_LOG = "log.txt"
const ARQUIVO_SITES = "sites.txt"

func main() {
	exibeIntroducao()

	for {
		exibeMenu()

		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			imprimeLogs()
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
	nome := "Weslley"
	versao := 0.1
	fmt.Println("Olá, ", nome, ". Este programa está na versão", versao)
}

func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do Programa")
}

func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)

	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando... ")

	sites := getSites()

	for i := 0; i < monitoramentos; i++ {
		for _, site := range sites {
			fmt.Println("Testando site:", site)
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
	fmt.Println("")
}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		fmt.Println("Ocorreu um erro: ", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registraLog(resp.StatusCode, site)
	} else {
		fmt.Println("Site:", site, "esta com problemas. Status Code:", resp.StatusCode)
		registraLog(resp.StatusCode, site)
	}

}

func getSites() []string {
	arquivo, err := os.Open(ARQUIVO_SITES)

	if err != nil {
		fmt.Println("Ocorreu um erro ao abrir o arquivo: ", err)
	}

	sites := leArquivo(arquivo)

	arquivo.Close()

	return sites
}

func leArquivo(arquivo *os.File) []string {
	var sites []string

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	return sites
}

func registraLog(statusCode int, site string) {
	arquivo, err := os.OpenFile(ARQUIVO_LOG, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Ocorreu um erro ao criar o arquivo: ", err)
	}

	date := time.Now().Format("02/01/2006 15:04:05")
	sc := strconv.Itoa(statusCode)

	arquivo.WriteString(date + " - " + sc + " - " + site + "\n")

	arquivo.Close()
}

func imprimeLogs() {

	arquivo, err := ioutil.ReadFile(ARQUIVO_LOG)

	if err != nil {
		fmt.Println("Ocorreu um erro ao ler o arquivo: ", err)
	}

	fmt.Println(string(arquivo))
}
