/* Definición de parser */


/* INICIA Sección de TOKENS */

%token a
%token b
/*FINALIZA Sección de TOKENS */

%%

/* INICIA Sección de PRODUCCIONES */

S:
    B C
  | D A
;

B:
    b
;

C:
    A A
;

A:
    a
;

D:
    b a
;

/* FINALIZA Sección de PRODUCCIONES
