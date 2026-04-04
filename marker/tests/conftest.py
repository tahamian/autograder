import sys
from pathlib import Path

# Add marker/ root to sys.path so `from marker import ...` works
sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
