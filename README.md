# CLI utility for automatic renaming pdf book files with isbn information from google books.

### How to use:

1. Copy to ```./data``` folder your pdf books.
2. Run commands in terminal:

    ```bash
    go build
    ./pdfbook-renamer -folder ./data
    ```


3. You PDFs in ```./data``` folder will be renamed to form: 

    ```title - author - publication year.pdf```

