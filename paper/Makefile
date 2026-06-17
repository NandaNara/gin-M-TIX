# ─────────────────────────────────────────────
# Makefile for gin-M-TIX IEEE Paper
# Requires: pdflatex + bibtex (e.g. TeX Live or MikTeX)
# ─────────────────────────────────────────────

MAIN    = main
LATEX   = pdflatex
BIBTEX  = bibtex
FLAGS   = -interaction=nonstopmode -halt-on-error

.PHONY: all clean view

all: $(MAIN).pdf

$(MAIN).pdf: $(MAIN).tex references.bib
	$(LATEX) $(FLAGS) $(MAIN)
	$(BIBTEX) $(MAIN)
	$(LATEX) $(FLAGS) $(MAIN)
	$(LATEX) $(FLAGS) $(MAIN)

clean:
	del /Q *.aux *.bbl *.blg *.log *.out *.toc *.lof *.lot 2>nul || true

view:
	start $(MAIN).pdf
