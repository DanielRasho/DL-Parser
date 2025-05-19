<h1 align="center">Parser and Lexer🚀</h1>
<h3 align="center">(A simple parser and lexer)</h3>

Hi! This project aims to be an educational implementation of a **Lexer Generator & Parser generator** following a Lex-like and Yapar-like sintax to configure it. 

It uses a **Deterministic Finite Automata(DFA)** for Regex Patterns recognition, and a **L0 Automata** for syntax recognition. Down below, there will be more explanation about the actual pipeline the inputs suffer in order to recognize tokens.

## Architecture videos
- [Part 1](https://youtu.be/WDLBOrFDXdQ)

## Getting Started 🎬

```bash
task buildParser    // Builds the parser generator
task test           // Run tests
task clean          // Removes executables
```

```bash
./bin/parser -f ./examples/productions3.y -o ./ora
# example:
./bin/parser -f ./examples/productions3.y -o ./ora
```

## The General Pipeline
A parser is a piece of software that identifies the structure of an input and tells you:

> "Dude, this code you gave me doesn't follow the syntax of your language :p"

When constructing a language, the parser is in charge of the **syntactic analysis of the code.** Acting like a guard, it stops invalid code from advancing to the next steps, saving time and computational resources.

Like lexers, parsers can share many components, making it easy to standardize them and build a Parser Generator. **This is exactly what our Yapar does!** The general flow consists of:

1. Providing a Yapar definition that specifies all the syntactic rules input code must follow to be considered valid.

2. A template that contains the common pieces all parsers share (check ours at `/template`).

3. The generator takes those rules and builds a functional `parser.go` file.

4. Providing a `lexer.go` (this one can also be generated) that reads the input code and pass it to the `parser.go` as a stream of tokens.

![](./pictures/parserPipeline.png)


## Parser Architecture


![](./pictures/parserArchitecture.png)

What is the parser doing?

1. Reads a Yapar file to recognize which tokens and productions can it recognize to either accept or not. 

2.  Those tokens and productions will pass them to the generator in order to create a table to check if the input we will give is either correct or not.  

    2.1 First step is to create a L0 automata which guides us from one node to another node based on the tokens production to find the finals nodes which will be used to parse. 

    2.2 We also compute the first and follow in order to generate a table to reduce or go to so se if it can be accepted or not. 

3. Now we introduce all componentes from the generator to a Parser template for GO in order to compile it and get a file to run it. 


## Data Structures

If you are more curious of how this was implemented on the code, you may start with the data types definitions. Almost every go module created has its own `types.go` file with the most important type definitions that modules exposes. Here are a list of the most importants:

### Yapar definition
https://github.com/DanielRasho/DL-Parser/blob/04793e148851f7b11137f49fbcca6fd51c9d85fc/internal/Parser/types.go#L9-L24

### Transition Table

https://github.com/DanielRasho/DL-Parser/blob/04793e148851f7b11137f49fbcca6fd51c9d85fc/internal/Parser/TransitionTable/types.go#L3-L25

### SLR0 Automata

https://github.com/DanielRasho/DL-Parser/blob/04793e148851f7b11137f49fbcca6fd51c9d85fc/internal/Parser/automata/types.go#L10-L23

### EXAMPLE OF ALL ARQUITECTURE
![image](https://github.com/user-attachments/assets/f75813fe-cb3a-46b8-a7d0-85f5e71a192b)
First, we create the YAPAR file where we identify the tokens and the productions that need to be parsed. We begin by reading the tokens and determining whether they are terminals. Then, there is a delimiter that signals when to start reading the productions. All of this is passed to a definition that identifies whether it is a non-terminal, terminal, and its corresponding productions.

![image](https://github.com/DanielRasho/DL-Parser/blob/main/pictures/First_Follow.png)
We then proceed with the First and Follow sets, where we find that First(E) is equal to First(T). After that, we calculate the Follow set in order to use it in the transition table.

![image](https://github.com/DanielRasho/DL-Parser/blob/main/pictures/Automata.png)
Before performing the transitions, we also need to generate the automaton, where we identify the accepting and final states in order to compute the shift and reduce actions.

![image](https://github.com/user-attachments/assets/dbf329d0-e132-4a36-a4e5-cb23d5669668)
With the transition table, we need the First and Follow sets, as well as the automaton, in order to compute the transition table and the GOTO table. These will be used to return a structure that should be embedded into a Go template for the purpose of parsing the input from a file.

![image](https://github.com/user-attachments/assets/b6460113-f9c8-4bc6-90f6-dd85c9ea5709)
So, once we integrate the transition table and the productions, we generate them in a main.go file that reads the input file, which contains the string to be parsed. Therefore, in our example, we use int + int, where we can see that the string must be accepted.


