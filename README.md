# Simple utility for automatic renaming pdf book files with isbn information from google books

### How to start:

1. Copy to ./data folder your pdf books.
2. Run command in terminal:

 ```docker build -t pdfbook-renamer . && docker run -it --rm -v "$(pwd)/data:/data" pdfbook-renamer:latest /app/pdfbook-renamer```

3. You pdf will be renamed to form: 

```title - author - publication year.pdf```

