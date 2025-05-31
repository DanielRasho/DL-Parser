// ======= HEADER =======
%{
    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file
    //  ------ TOKENS ID -----
    // Define the token types that the lexer will recognize
    const (
        LET = iota
        ASSIGN
        PLUS
        MINUS
        MULT
        DIV
        ID
        NUMBER
        SEMICOLON
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    digit        [0-2]
    letter       [a-c]
    id           {letter}({letter}|{digit})*
    number       ({digit})+
    ws           ([ \t\n\r])+
}

// ======= RULES ========
%%
"let"            { return LET }
"="             { return ASSIGN }
"\+"           { return PLUS }
"-"             { return MINUS }
"\*"           { return MULT }
"/"             { return DIV }
";"             { return SEMICOLON }

{id}            { return ID }
{number}        { return NUMBER }
{ws}            {} 
%%

// ======= FOOTER =======
%{
    // This is a footer section where additional methods can be added if needed.
%}