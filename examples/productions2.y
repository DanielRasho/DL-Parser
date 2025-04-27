/* Definición de parser */


/* INICIA Sección de TOKENS */

%token int
%token +
%token *
%token (
%token )
/*FINALIZA Sección de TOKENS */

%%

/* INICIA Sección de PRODUCCIONES */
E:
    T + E
  | T
;

T:
    int * T
  | int
  | ( E )
;
/* FINALIZA Sección de PRODUCCIONES
