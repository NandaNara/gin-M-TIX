# gin-M-TIX — IEEE Paper

This directory contains a full-length **IEEE conference paper** written in LaTeX about the three design patterns implemented in the `gin-M-TIX` Cinema Ticket Booking API.

## Files

| File | Description |
|------|-------------|
| `main.tex` | Full IEEE-format LaTeX source paper |
| `references.bib` | BibTeX bibliography (7 references) |
| `Makefile` | Build script (Windows `pdflatex` + `bibtex`) |

## Paper Overview

**Title:** *Implementation of Factory, Strategy, and Facade Design Patterns in a Go-Based Cinema Ticket Booking REST API*

**Abstract:** Demonstrates how three GoF patterns (Factory Method, Strategy, Facade) are applied in idiomatic Go within a real REST API, with annotated code listings, UML structure descriptions, and an extensibility evaluation.

### Sections

1. Introduction
2. Background & Related Work
3. System Architecture
4. **Factory Method Pattern** — `RegularTicketFactory`, `VIPTicketFactory`, `StudentTicketFactory`
5. **Strategy Pattern** — `HolidayPricing`, `MidnightPricing`, `WeekendPricing`, `WeekdayPricing`
6. **Facade Pattern** — `BookingFacade` over `BookingService` + `PaymentService`
7. Evaluation (OCP compliance, coupling analysis, Go idioms, limitations)
8. Conclusion

## How to Compile

### Prerequisites
Install one of:
- [MikTeX](https://miktex.org/) (Windows — recommended, auto-installs packages)
- [TeX Live](https://www.tug.org/texlive/)

### Build

```bash
# Using Makefile (requires make — e.g. via Git Bash or Chocolatey)
make

# Or manually with PowerShell / cmd:
pdflatex -interaction=nonstopmode main
bibtex main
pdflatex -interaction=nonstopmode main
pdflatex -interaction=nonstopmode main
```

The output will be `main.pdf`.

### Open PDF
```bash
make view
# or
start main.pdf
```

## Customization

- Replace `[Your Name]` and `[Your University]` in `main.tex` (lines under `\author{}`) with your actual details.
- Add UML diagrams by exporting them as `.pdf` or `.png` and referencing with `\includegraphics{}`.
- The paper is structured for an **IEEE conference** (`\documentclass[conference]{IEEEtran}`). For a journal submission, change to `\documentclass[journal]{IEEEtran}`.
