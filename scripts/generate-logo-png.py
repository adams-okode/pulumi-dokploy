#!/usr/bin/env python3
"""Render the repository's small, static logo SVG without external tools."""

import re
import struct
import sys
import zlib
from pathlib import Path
from xml.etree import ElementTree

WIDTH = HEIGHT = 175
VIEWBOX = (0.0, 0.0, 457.0, 472.0)
SAMPLES = 4
NUMBER = r"[-+]?(?:\d+(?:\.\d*)?|\.\d+)(?:[eE][-+]?\d+)?"


def parse_logo(path: Path) -> tuple[list[tuple[float, float]], tuple[int, int, int]]:
    root = ElementTree.fromstring(path.read_bytes())
    if root.tag.rsplit("}", 1)[-1] != "svg" or root.attrib.get("viewBox") != "0 0 457 472":
        raise ValueError("logo SVG must use the canonical 457 by 472 viewBox")
    paths = [node for node in root.iter() if node.tag.rsplit("}", 1)[-1] == "path"]
    if len(paths) != 1 or paths[0].attrib.keys() != {"fill", "d"}:
        raise ValueError("logo SVG must contain exactly one plain filled path")
    color = paths[0].attrib["fill"]
    if not re.fullmatch(r"#[0-9a-fA-F]{6}", color):
        raise ValueError("logo path fill must be an RGB hex color")
    tokens = re.findall(rf"[A-Za-z]|{NUMBER}", paths[0].attrib["d"])
    if not tokens or tokens[0] != "M" or tokens[-1] != "z":
        raise ValueError("logo path must start with move and end with close")
    if any(token not in {"M", "c", "l", "z"} for token in tokens if token.isalpha()):
        raise ValueError("logo path contains an unsupported command")
    cursor = 0
    points = []
    command = None
    x = y = 0.0
    while cursor < len(tokens):
        if tokens[cursor].isalpha():
            command = tokens[cursor]
            cursor += 1
        if command == "z":
            break
        if command == "M":
            x, y = float(tokens[cursor]), float(tokens[cursor + 1])
            cursor += 2
            points.append((x, y))
            command = "l"
        elif command == "l":
            x += float(tokens[cursor])
            y += float(tokens[cursor + 1])
            cursor += 2
            points.append((x, y))
        elif command == "c":
            controls = [float(value) for value in tokens[cursor : cursor + 6]]
            if len(controls) != 6:
                raise ValueError("logo path has incomplete cubic geometry")
            cursor += 6
            c1 = (x + controls[0], y + controls[1])
            c2 = (x + controls[2], y + controls[3])
            end = (x + controls[4], y + controls[5])
            start = (x, y)
            for step in range(1, 17):
                t = step / 16
                inverse = 1 - t
                points.append(tuple(
                    inverse**3 * start[index]
                    + 3 * inverse**2 * t * c1[index]
                    + 3 * inverse * t**2 * c2[index]
                    + t**3 * end[index]
                    for index in (0, 1)
                ))
            x, y = end
        else:
            raise ValueError("logo path has invalid geometry")
    return points, tuple(int(color[index : index + 2], 16) for index in (1, 3, 5))


def inside(point: tuple[float, float], polygon: list[tuple[float, float]]) -> bool:
    x, y = point
    result = False
    for index, (x1, y1) in enumerate(polygon):
        x2, y2 = polygon[index - 1]
        if (y1 > y) != (y2 > y) and x < (x2 - x1) * (y - y1) / (y2 - y1) + x1:
            result = not result
    return result


def png_chunk(kind: bytes, data: bytes) -> bytes:
    return struct.pack(">I", len(data)) + kind + data + struct.pack(">I", zlib.crc32(kind + data) & 0xFFFFFFFF)


def render(points: list[tuple[float, float]], color: tuple[int, int, int]) -> bytes:
    scale_x = WIDTH / VIEWBOX[2]
    scale_y = HEIGHT / VIEWBOX[3]
    rows = []
    for y in range(HEIGHT):
        row = bytearray([0])
        for x in range(WIDTH):
            covered = 0
            for sy in range(SAMPLES):
                for sx in range(SAMPLES):
                    source = ((x + (sx + 0.5) / SAMPLES) / scale_x,
                              (y + (sy + 0.5) / SAMPLES) / scale_y)
                    covered += inside(source, points)
            alpha = round(255 * covered / (SAMPLES * SAMPLES))
            row.extend((*color, alpha))
        rows.append(row)
    header = struct.pack(">IIBBBBB", WIDTH, HEIGHT, 8, 6, 0, 0, 0)
    return b"\x89PNG\r\n\x1a\n" + png_chunk(b"IHDR", header) + png_chunk(b"IDAT", zlib.compress(b"".join(rows), 9)) + png_chunk(b"IEND", b"")


def main() -> int:
    if len(sys.argv) != 3:
        raise SystemExit("usage: generate-logo-png.py SVG OUTPUT")
    points, color = parse_logo(Path(sys.argv[1]))
    output = Path(sys.argv[2])
    output.parent.mkdir(parents=True, exist_ok=True)
    temporary = output.with_suffix(output.suffix + ".tmp")
    temporary.write_bytes(render(points, color))
    temporary.replace(output)
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, ElementTree.ParseError, ValueError) as error:
        raise SystemExit(f"logo rendering failed: {error}")
