"""Validate local asset and page references in the static portfolio."""
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[1]
EXCLUDED = {".git", ".github", "node_modules"}

class References(HTMLParser):
    def __init__(self):
        super().__init__()
        self.refs = []

    def handle_starttag(self, tag, attrs):
        for attr, value in attrs:
            if attr in {"href", "src"} and value:
                self.refs.append(value)

errors = []
scanned = 0
for page in sorted(ROOT.rglob("*.html")):
    if any(part in EXCLUDED for part in page.relative_to(ROOT).parts):
        continue
    scanned += 1
    parser = References()
    parser.feed(page.read_text(encoding="utf-8"))
    for ref in parser.refs:
        uri = urlsplit(ref)
        if uri.scheme or uri.netloc or ref.startswith(("#", "//")):
            continue
        if not uri.path:
            continue
        path = unquote(uri.path)
        target = ROOT / path.lstrip("/") if path.startswith("/") else page.parent / path
        if not target.resolve().is_relative_to(ROOT.resolve()) or not target.exists():
            errors.append(f"{page.relative_to(ROOT)}: missing {ref}")

if errors:
    raise SystemExit("Broken references:\n" + "\n".join(errors))
print(f"Checked {scanned} HTML files: all local links and assets exist")
