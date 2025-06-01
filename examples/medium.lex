// ======= HEADER =======
%{
    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file
    //  ------ TOKENS ID -----
    // Define the token types that the lexer will recognize
    const (
        LET = iota
        IF
        WHILE
        ASSIGN
        PLUS
        MINUS
        MULT
        DIV
        GT
        LT
        EQ
        LPAREN
        RPAREN
        LBRACE
        RBRACE
        ID
        NUMBER
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    digit        ([0-3])
    letter       ([a-d])
    id           {letter}({letter}|{digit})*
    number       ({digit})+
    WS           ([ \t\n\r])+
}

// ======= RULES ========
%%
"let"            { return LET }
"if"             { return IF }
"while"          { return WHILE }

"="             { return ASSIGN }
"\+"           { return PLUS }
"-"             { return MINUS }
"\*"           { return MULT }
"/"             { return DIV }

">"             { return GT }
"<"             { return LT }
"=="            { return EQ }

"\("           { return LPAREN }
"\)"           { return RPAREN }
"\{"           { return LBRACE }
"\}"           { return RBRACE }

{id}            { return ID }
{number}        { return NUMBER }
{WS}            { return WS } 
%%

// ======= FOOTER =======
%{
    // This is a footer section where additional methods can be added if needed.
%}