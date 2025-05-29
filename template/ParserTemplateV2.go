package main

func newTransitTable() *TransitionTbl {
	{{ .TransitTable }}

}

func newGoToTable() *GotoTbl {
	{{ .Gotable }}
}

