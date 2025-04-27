/* Definición de parser */


/* INICIA Sección de TOKENS */

%token ^
%token v
%token [
%token ]
%token sentence
/*FINALIZA Sección de TOKENS */

%%

/* INICIA Sección de PRODUCCIONES */

S:

    S ^ P
  | P
;

P:
    P v Q
  | Q
;

Q:
    [ S ]
  | sentence
;

/* FINALIZA Sección de PRODUCCIONES
