# pdfbook-renamer 
CLI utility for automatic renaming pdf book with ISBN info from google books.

## Prerequisites

1. Install Go:
   - Visit the official Go website: https://golang.org/
   - Download and install the appropriate version for your operating system
   - Verify the installation by running `go version` in your terminal

## How to use:

1. Run the following command to install the utility: 
    ```bash
    go install github.com/dmikoss/pdfbook-renamer
    ``` 
2. Copy to ```./folder-with-pdfs``` folder your pdf books.
3. Run commands in terminal:
    ```bash
    pdfbook-renamer -folder ./folderwithpdfs
    ```
4. You PDFs in ```./folder-with-pdfs``` folder will be renamed to form: 

    ```title - author - publication_year.pdf```

