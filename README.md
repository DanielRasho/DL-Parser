<h1 align="center">Parser and Lexer🚀</h1>
<h3 align="center">(A simple parser and lexer)</h3>

Hi! This project aims to be an educational implementation of a **Lexer Generator & Parser generator** following a Lex-like and Yapar-like sintax to configure it. 

Please refer to the Paser or Lexer folders of the repo to a more in-depth explanation of its functionality

- [**Lexer**](https://github.com/DanielRasho/DL-Parser/tree/main/internal/Lexer)
- [**Parser**](https://github.com/DanielRasho/DL-Parser/tree/main/internal/Parser)

## Getting Started 🎬
You must have go-task, go, and graphviz on your computer. If you are using Linux, Mac or WSL, this project is configured with Nix, to download the dependencies mencioned above for you. Just run:

```bash
nix develop
```
And you are ready to rock!

For the build commands, we are using go-task, to show available commands run
```bash
task
```
An example of command look like this:
```bash
task compile:build -- -y yapar.y -l yalex.l -c source.txt -o ./compiler
```
