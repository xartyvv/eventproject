import pathlib
import re

root = pathlib.Path('.')
emoji_pattern = re.compile('[\U0001F300-\U0001F6FF\U0001F900-\U0001F9FF\u2600-\u26FF\u2700-\u27BF]')
for path in root.rglob('*'):
    if path.is_file() and path.suffix.lower() in ['.html', '.js', '.css', '.md', '.bat']:
        text = path.read_text(encoding='utf-8', errors='ignore')
        if emoji_pattern.search(text):
            print(path, emoji_pattern.findall(text)[:20])
