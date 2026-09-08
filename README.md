# Language Refresher

Aplicação desktop para sessões curtas e práticas de revisão de linguagens de programação. A primeira trilha é **Python Core Refresher**, com 45 exercícios e duração estimada de aproximadamente duas horas.

## Stack

- Go + Wails v2 no backend desktop
- React + TypeScript no frontend
- CodeMirror 6 no editor
- Conteúdo e progresso em JSON

O Wails usa o WebView nativo do sistema, entrega uma aplicação desktop real sem Electron e permite manter execução, timeout, cancelamento e persistência fora da camada de UI. O CodeMirror fornece highlighting, números de linha, indentação, histórico e atalhos com uma integração pequena.

## Desenvolvimento

Pré-requisitos: Go 1.23+, Node.js/npm, Python 3.10+ e WebView2 no Windows.

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@v2.15.0
wails doctor
wails dev
```

Para gerar o executável:

```powershell
wails build
```

Atalhos principais:

- `Ctrl+Enter`: executar
- `Ctrl+S`: salvar o rascunho
- atalhos nativos do editor: selecionar tudo, copiar, colar, desfazer e refazer
- arraste o divisor vertical para redimensionar o enunciado
- arraste o divisor horizontal para redimensionar editor e terminal
- use as setas quando um divisor estiver focado; clique duas vezes para restaurar o tamanho padrão

## Idiomas

Na primeira abertura, a interface usa português para localidades `pt-*` e inglês para qualquer outra localidade. A engrenagem ao lado do nome da aplicação permite alternar entre **Português** e **English**; a escolha é persistida em `settings.json` no diretório de configuração do usuário.

As traduções dos exercícios ficam em arquivos `*.i18n.json`, separadas de starter code, arquivos de apoio e validadores. Dessa forma, trocar o idioma preserva o exercício atual, o rascunho e todo o progresso.

## Arquitetura

```text
internal/
  exercises/   carregamento e validação do catálogo
  execution/   processos, captura, timeout e cancelamento
  languages/   comandos e arquivos por linguagem
  progress/    persistência JSON local
  sessions/    estado e navegação da trilha
content/
  python/      trilhas e testes declarativos
frontend/      interface React/CodeMirror
```

Os testes automáticos executam em um diretório temporário separado, carregam `main.py` por `runpy` e inspecionam funções, classes e valores. Isso evita depender apenas da comparação de stdout. O MVP **não é uma sandbox de segurança**: o código continua com as permissões do usuário local. A camada `execution` foi isolada para permitir substituir o executor por contêiner, VM ou outro isolamento no futuro.
