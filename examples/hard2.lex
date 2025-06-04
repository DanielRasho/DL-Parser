// ======= HEADER =======
%{
    const (
        IF = iota
        ELSE
        WHILE
        RETURN
        FUNC
        VAR
        INT
        FLOAT
        STRING
        BOOL
        TRUE
        FALSE
        FOR
        BREAK
        CONTINUE
        ASSIGN
        PLUS
        MINUS
        MULT
        DIV
        MOD
        AND
        OR
        NOT
        EQ
        NEQ
        LT
        GT
        LTE
        GTE
        LPAREN
        RPAREN
        LBRACE
        RBRACE
        LBRACKET
        RBRACKET
        COMMA
        SEMICOLON
        DOT
        ID
        NUMBER
        FLOAT_LIT
        STRING_LIT
        COMMENT
        WS
    )
%}

// ====== NAMED PATTERNS =======
{
    digit         ([0-9])
    letter        ([a-zA-Z_])
    id            {letter}({letter}|{digit})*
    number        ({digit})+
    float_lit     ({digit})+.({digit})+
    string_lit    "(\\.|[^\\"])*"
    single_comment    (\/\/.)*
    WS            ([ \t\n\r])+
}

// ======= RULES ========
%%
"if"                { return IF }
"else"              { return ELSE }
"while"             { return WHILE }
"return"            { return RETURN }
"func"              { return FUNC }
"var"               { return VAR }
"int"               { return INT }
"float"             { return FLOAT }
"string"            { return STRING }
"bool"              { return BOOL }
"true"              { return TRUE }
"false"             { return FALSE }
"for"               { return FOR }
"break"             { return BREAK }
"continue"          { return CONTINUE }

"="                 { return ASSIGN }
"\+"                { return PLUS }
"-"                 { return MINUS }
"\*"                { return MULT }
"/"                 { return DIV }
"%"                 { return MOD }
"&&"                { return AND }
"\|\|"              { return OR }
"!"                 { return NOT }
"=="                { return EQ }
"!="                { return NEQ }
"<"                 { return LT }
">"                 { return GT }
"<="                { return LTE }
">="                { return GTE }

"\("                { return LPAREN }
"\)"                { return RPAREN }
"\{"                { return LBRACE }
"\}"                { return RBRACE }
"\["                { return LBRACKET }
"\]"                { return RBRACKET }
","                 { return COMMA }
";"                 { return SEMICOLON }
"\."                { return DOT }

{float_lit}         { return FLOAT_LIT }
{number}            { return NUMBER }
{string_lit}        { return STRING_LIT }
{single_comment}    { return COMMENT }
{multi_comment}     { return COMMENT }
{id}                { return ID }
{WS}                { return WS }
%%

// ======= FOOTER =======
%{
    // Define helper functions or add error-handling logic here
%}