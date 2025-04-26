/* Definición de parser */


/* INICIA Sección de TOKENS */

%token IF ELSE THEN END
%token IDENTIFIER
%token WS 
IGNORE WS
/* FINALIZA Sección de TOKENS */


%%

/* INICIA Sección de PRODUCCIONES */

statement:
    IF condition THEN block ELSE block END
  | IF condition THEN block END
;

condition:
    IDENTIFIER
;

block:
    statement
  | IDENTIFIER
;

/* FINALIZA Sección de PRODUCCIONES */
