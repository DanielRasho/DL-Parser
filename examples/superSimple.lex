// ======= HEADER =======
%{
    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file
    //  ------ TOKENS ID -----
    // Define the token types that the lexer will recognize
    const (
        DIGIT = iota
        PLUS
        MULT
        LPAREN
        RPAREN
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    DIGIT   int
    PLUS    \+
    MULT    \*
    LPAREN  \(
    RPAREN  \)
    WS  ([ \t\n\r])+
}

// ======= RULES ========
%%
{DIGIT}         { return DIGIT }   // Match letters and return LITERAL
{PLUS}          { return PLUS }    // Match digits and return NUMBER
{MULT}          { return MULT }
{LPAREN}        { return LPAREN }
{RPAREN}        { return RPAREN }
{WS}            { return WS }
%%
