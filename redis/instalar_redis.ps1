# Instalar WSL e Ubuntu no Windows

Write-Host "Instalando WSL e a distribuição Ubuntu..."
wsl --install -d Ubuntu

Write-Host "Reinicie o computador após a instalação do WSL."
Pause

# Mensagem de instrução após a reinicialização
Write-Host "Após a reinicialização, abra o Ubuntu e siga os passos na tela para configurar."

Pause
Write-Host "Agora vamos instalar o Redis no Ubuntu."

# Comando para abrir o terminal do Ubuntu no WSL e instalar o Redis
wsl -d Ubuntu sudo apt update
wsl -d Ubuntu sudo apt install -y redis-server

Write-Host "Redis instalado. Agora vamos iniciar o servidor Redis."

# Iniciar o servidor Redis
wsl -d Ubuntu sudo service redis-server start

Write-Host "Servidor Redis iniciado."

# Testar se o Redis está funcionando
Write-Host "Testando a conexão com o Redis..."
$redisPing = wsl -d Ubuntu redis-cli ping

if ($redisPing -eq "PONG") {
    Write-Host "Redis está funcionando corretamente!"
} else {
    Write-Host "Algo deu errado. Redis não respondeu com PONG."
}

Pause

Write-Host "Agora vamos testar o cache do Redis."

# Adicionar e obter valores do cache Redis
wsl -d Ubuntu redis-cli set minhaChave "Este é um valor no cache"
$cacheValue = wsl -d Ubuntu redis-cli get minhaChave

if ($cacheValue -eq "Este é um valor no cache") {
    Write-Host "Cache testado com sucesso! Valor armazenado: $cacheValue"
} else {
    Write-Host "Erro ao testar o cache."
}

Write-Host "Teste completo. O Redis foi instalado e o cache está funcionando."

Pause

Write-Host "Se desejar, você pode parar o Redis com o seguinte comando:"
Write-Host "sudo service redis-server stop"
