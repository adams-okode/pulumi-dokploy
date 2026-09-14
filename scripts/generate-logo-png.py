#!/usr/bin/env python3
"""Materialize the checked-in, byte-stable Registry logo PNG.

The Registry accepts a PNG while the website owns the SVG source.  Keeping the
renderer output under version control avoids depending on a distro package
index whose contents vary with the GitHub runner image.
"""

import hashlib
import subprocess
import sys
from pathlib import Path

SVG_SHA256 = "b87290aa2cc570aec18e4f0d640bddd37fc6edbfba5a2be503bf11ce81dfc375"


def main() -> int:
    if len(sys.argv) != 3:
        raise SystemExit("usage: generate-logo-png.py SVG OUTPUT")
    svg, output = map(Path, sys.argv[1:])
    if hashlib.sha256(svg.read_bytes()).hexdigest() != SVG_SHA256:
        raise SystemExit(f"{svg} changed; regenerate and review the canonical logo PNG")
    png = subprocess.run(
        ["git", "show", "HEAD:sdk/dotnet/logo.png"],
        check=True,
        capture_output=True,
    ).stdout
    output.parent.mkdir(parents=True, exist_ok=True)
    temporary = output.with_suffix(output.suffix + ".tmp")
    temporary.write_bytes(png)
    temporary.replace(output)
    return 0


if __name__ == "__main__":
    sys.exit(main())
