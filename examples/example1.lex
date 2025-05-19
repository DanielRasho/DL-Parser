// ======= HEADER =======
%{
    // The entire contents of this section will be copied to the beginning of the generated Lexer.go file
    //  ------ TOKENS ID -----
    // Define the token types that the lexer will recognize
    const (
        LITERAL = iota
        NUMBER
        DIGIT
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    // Define named patterns using regular expressions
    LETTER   ([a-b])+
    DIGIT    ([1-2])
    NUMBER   {DIGIT}+  // NUMBER consists of one or more digits
    WS       ([ \t\n])+  // Whitespace: spaces, tabs, newlines, or carriage returns
}

// ======= RULES ========
%%
{LETTER}         { return LITERAL }   // Match letters and return LITERAL
{NUMBER}          { return NUMBER }    // Match digits and return NUMBER
{WS}         { return WS }

%%
