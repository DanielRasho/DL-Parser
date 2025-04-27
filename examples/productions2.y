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

T:
    int * T
  | int
  | ( E )
;

E:

    T + E
  | T
;
/* FINALIZA Sección de PRODUCCIONES
