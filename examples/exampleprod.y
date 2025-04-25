/* Definición de parser */


/* INICIA Sección de TOKENS */

%token NUMBER
%token PLUS MINUS TIMES DIVIDE
%token LPAREN RPAREN
%token WS expr term factor
IGNORE WS
/* FINALIZA Sección de TOKENS */


%%

/* INICIA Sección de PRODUCCIONES */

expr:
    expr PLUS term
  | expr MINUS term
  | term
;

term:
    term TIMES factor
  | term DIVIDE factor
  | factor
;

factor:
    LPAREN expr RPAREN
  | NUMBER
;

/* FINALIZA Sección de PRODUCCIONES */
