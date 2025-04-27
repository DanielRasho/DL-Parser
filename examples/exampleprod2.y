/* Definición de parser */


/* INICIA Sección de TOKENS */

%token IF ELSE THEN END
%token IDENTIFIER
%token WS 
IGNORE WS
/* FINALIZA Sección de TOKENS */


%%

/* INICIA Sección de PRODUCCIONES */

S:
    IF C THEN B ELSE B END
  | IF C THEN B END
;

C:
    IDENTIFIER
;

B:
    S
  | IDENTIFIER
;

/* FINALIZA Sección de PRODUCCIONES */
