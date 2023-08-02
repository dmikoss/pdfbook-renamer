import sys
from pypdf import PdfReader

if len(sys.argv) < 3:
    sys.exit(-2)

pdffile = sys.argv[1]
reader = PdfReader(pdffile)
number_of_pages = len(reader.pages)
num_pages_to_scan = int(sys.argv[2])

for i in range(num_pages_to_scan):
    page = reader.pages[i]
    text = page.extract_text()
    print (text)
