/* Definición de parser */


/* INICIA Sección de TOKENS */

%token TOKEN_1
%token TOKEN_2
%token TOKEN_3 TOKEN_4
%token WS 
IGNORE WS
/*FINALIZA Sección de TOKENS */



%%

/* INICIA Sección de PRODUCCIONES */
production1:
    production1 TOKEN_2 production2 
  | production2
;

production2:
    production2 TOKEN_2 production3 
  | production3
;


production3:
    TOKEN_3 production1 TOKEN_4 
  | TOKEN_1
;
/* FINALIZA Sección de PRODUCCIONES
