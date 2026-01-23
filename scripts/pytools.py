import os

import docx
import pdfplumber
import uvicorn
import requests

from fastapi import FastAPI

app = FastAPI()


@app.get("/extract-file-text")
def extract(url: str):
    filename = url.split('/')[-1]
    filetype = filename.split('.')[-1]

    try:
        response = requests.get(url, verify=False)
        with open(filename, 'wb') as file:
            file.write(response.content)

        if filetype == 'pdf':
            with pdfplumber.open(filename) as pdf:
                first_page = pdf.pages[0]
                return {'data': first_page.extract_text()}
        if filetype == 'docx' or filetype == 'doc':

            document = docx.Document(filename)
            text = ''
            all_paragraphs = document.paragraphs
            for paragraph in all_paragraphs:
                text += paragraph.text + '\n\n'

            return {'data': text}
    finally:
        if os.path.exists(filename):
            os.remove(filename)


if __name__ == "__main__":
    uvicorn.run(app="pytools:app", host="0.0.0.0", port=8080)
