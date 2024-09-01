package main

import (
	// Pacotes nativos
	"bytes"         // Fornece funções para manipular e gerenciar slices de bytes (sequências de bytes). Ele é útil, por exemplo, para criar e manipular o corpo de uma requisição HTTP em Go.
	"context"       // Permite a criação de contextos, que são usados para cancelar operações, propagar informações de requisições e gerenciar o tempo de execução.
	"encoding/json" // Responsável pela codificação e decodificação de dados em formato JSON. Ele permite converter estruturas de dados do Go para JSON e vice-versa.
	"fmt"           // Frequentemente usado para imprimir dados no console e formatar strings.
	"io/ioutil"     // Fornece funções para trabalhar com entrada e saída de dados (I/O). Ele permite ler e gravar dados de arquivos, streams, etc
	"log"           // Registra mensagens de log, que são úteis para depuração e monitoramento de aplicações.
	"net/http"      // Ele lida com o cliente HTTP, requisições, respostas e servidores.
	"os"            // Manipulação de arquivos, variáveis de ambiente e processos.
	"sync"          // Fornece mecanismos de sincronização para concorrência, como mutexes e grupos de espera.
	"time"          // Ele permite pausar a execução, medir durações e definir TTL (tempo de vida) de cache.

	// Pacotes de terceiros
	"github.com/gin-gonic/gin"     // Framework web escrito em Go (Golang)
	"github.com/joho/godotenv"     // Usado para carregar variáveis de ambiente a partir de arquivos .env
	"github.com/redis/go-redis/v9" // Cliente Redis para a linguagem Go. Ele permite que você interaja com um banco de dados Redis diretamente a partir de suas aplicações Go. Este cliente é mantido pela comunidade e é amplamente utilizado devido à sua eficiência e facilidade de uso.
)

var (
    ctx    = context.Background() // Espaço vazio criado. Será utilizado para armazenar informações do redis (https://pkg.go.dev/context#Background)
    client *redis.Client // Variável que representa uma conexão com o banco de dados Redis
    mutex  sync.Mutex // Mutex global para garantir o comportamento síncrono
    isLocked   bool // Variável auxiliar para verificar se o mutex está bloqueado
)

func init() {
    // Carrega as variáveis de ambiente do arquivo .env
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
    }

    // Configura o cliente Redis
    client = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"), // Endereço do Redis  os.Getenv("REDIS_ADDR")
        Password: os.Getenv("REDIS_PASSWORD"), // Sem senha os.Getenv("REDIS_PASSWORD")
        DB:       0,                // Banco de dados padrão
    })
}

// Estrutura para o corpo da requisição
type OmieRequest struct {
    Call      string                 json:"call";
    Param     []map[string]interface{} json:"param";
    URL       string                 json:"url" // Adicionamos o campo URL
}

// Função para fazer a requisição à API OMIE e retornar a resposta bruta
func fetchFromOmie(requestBody OmieRequest) ([]byte, error) {
    baseURL := "https://app.omie.com.br/api/v1/"
    url := baseURL + requestBody.URL

    // Obter as chaves OMIE das variáveis de ambiente
    appKey := os.Getenv("OMIE_APP_KEY")
    appSecret := os.Getenv("OMIE_APP_SECRET")

    // Adiciona as chaves ao request body
    requestBodyWithAuth := map[string]interface{}{
        "call":       requestBody.Call,
        "app_key":    appKey,
        "app_secret": appSecret,
        "param":      requestBody.Param,
    }

    body, err := json.Marshal(requestBodyWithAuth) // Converte a estrutura go em byte slice contendo JSON. 
    if err != nil {
        return nil, err
    }

    // Cria um contexto com timeout de 1 minuto
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close() // Fecha o corpo da resposta 'resp.Body' após a leitura dos dados.

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API retornou código de status %d", resp.StatusCode)
    }

    // Lê a resposta da API OMIE como um JSON bruto
    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return responseBody, nil
}

// Função para obter dados com cache
func getData(requestBody OmieRequest) ([]byte, error) {
    // Cria uma chave única para o cache baseada na requisição
    cacheKey, err := json.Marshal(requestBody) // Converte a estrutura go em byte slice contendo JSON. 
    if err != nil {
        return nil, err
    }

    //Para testar o delay, caso existir
    // time.Sleep(2 * time.Second)

    // Verifica o cache. 'client.Get' realiza leitura dos dados do caché
    cachedData, err := client.Get(ctx, string(cacheKey)).Result() // 'Result()' devolve um valor do redis
    if err == nil {
        fmt.Println("Dados recuperados do cache.")
        // Retorna os dados do cache diretamente como JSON bruto
        return []byte(cachedData), nil
    }

    // Dados não encontrados no cache, solicita à API OMIE
    data, err := fetchFromOmie(requestBody)
    if err != nil {
        return nil, err
    }

    // Armazena no cache por 60 segundos. 'client.Set' realiza a gravação de dados
    err = client.Set(ctx, string(cacheKey), string(data), 60*time.Second).Err()
    if err != nil {
        return nil, err
    }

    fmt.Println("Dados recuperados diretamente da API Omie.")
    return data, nil
}

func main() {
    router := gin.Default() //  Inicializa o roteador HTTP da biblioteca Gin
    fmt.Println("↓↓↓↓↓↓ ▲▼ Damas e Cavalheiros, o servidor esta rodando ▲▼ ↓↓↓↓↓↓")
    router.POST("/omie_request", func(c *gin.Context) { // Definida uma rota HTTP do tipo POST associada ao caminho /omie_request
        var requestBody OmieRequest
        if err := c.BindJSON(&requestBody); err != nil { // 'BindJSON'  tenta converter (ou fazer "unmarshal") o corpo da requisição JSON em uma estrutura Go. 'c' é uma instância ou espaço para mexer com dados
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"}) // Se a estrutura esta correta, prosseguirá, se esta incorreta, bloqueara a requisição
            return
        }

        // Verifica o estado do mutex
        if isLocked {
            fmt.Println("Processando requisição. Aguardando liberações anteriores...")
        }

        // Bloqueia o mutex. A nova requisição será forçada a esperar até que o mutex seja liberado. Todo o código que envolve verificar o cache, consultar a API Omie, armazenar no cache e retornar os dados está protegido pelo mutex
        mutex.Lock()
        isLocked = true // Marca que o mutex está bloqueado

        // Garante que o mutex será liberado e que o estado será restaurado ao final
        defer func() {
            // Após todo o processo (cache, API, armazenar e retornar), o mutex é liberado, permitindo que a próxima requisição na fila entre e comece o processo.
            mutex.Unlock()
            isLocked = false // Marca que o mutex foi liberado
            fmt.Println("Requisição processada. Mutex liberado.")
        }()


        // Busca os dados: verifica o cache, faz requisição e volta com um response
        data, err := getData(requestBody)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        fmt.Println("Dados recebidos:", string(data))
        // Retorna a resposta ao cliente (Python) como JSON bruto
        c.Data(http.StatusOK, "application/json", data)
    })

    router.Run(":8080") // O servidor vai rodar na porta 8080
}

// https://redis.io/docs/latest/commands/expire/ DOCUMENTAÇÃO DE COMO FUNCIONA O TIMESTAMPS
// https://redis.io/docs/latest/operate/oss_and_stack/management/optimization/memory-optimization/ MEMORIA OTIMIZADA
// https://pkg.go.dev/time ENTENDER CRONOMETRO - Cria um timStamp para colocar em uma chave no Redis