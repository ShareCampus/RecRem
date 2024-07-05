import fitz  # PyMuPDF
import docx
import sys

def extract_text_from_pdf(pdf_path):
    document = fitz.open(pdf_path)
    text = ""
    for page_num in range(len(document)):
        page = document.load_page(page_num)
        text += page.get_text()
    return text

def extract_text_from_docx(docx_path):
    doc = docx.Document(docx_path)
    text = ""
    for paragraph in doc.paragraphs:
        text += paragraph.text + "\n"
    return text

def extract_text_from_txt(txt_path):
    with open(txt_path, 'r', encoding='utf-8') as file:
        text = file.read()
    return text

def extract_text_from_md(md_path):
    with open(md_path, 'r', encoding='utf-8') as file:
        text = file.read()
    return text

def extract_text(file_path):
    if file_path.endswith('.pdf'):
        return extract_text_from_pdf(file_path)
    elif file_path.endswith('.docx'):
        return extract_text_from_docx(file_path)
    elif file_path.endswith('.txt'):
        return extract_text_from_txt(file_path)
    elif file_path.endswith('.md'):
        return extract_text_from_md(file_path)
    else:
        raise ValueError("Unsupported file format: only .pdf, .docx, .txt, and .md are supported")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python extract_text.py <file_path>")
        sys.exit(1)

    file_path = sys.argv[1]
    try:
        text = extract_text(file_path)
        print(text)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)