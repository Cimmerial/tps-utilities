# tps-utilities

General utilities run through [TPS](https://github.com/Cimmerial/tps). **Do not clone this repo alone expecting it to work** — install TPS first, or run `./setup` after clone (it checks for `tps` on PATH).

## Install from TPS

In any `tps` session:

```text
iu
```

(`install utilities` — picks a parent folder on macOS, clones this repo, links the project.)

Switch to the utilities project (`p`) and run **`uic`** (Universal Image Conversion).

## Direct clone

```bash
git clone https://github.com/Cimmerial/tps-utilities.git
cd tps-utilities
chmod +x setup
./setup
```

Requires **ImageMagick** (`magick`). On first convert, `tpsu` can install it for you (Homebrew on macOS, apt/dnf/pacman on Linux, winget/choco/scoop on Windows). Non-interactive: `TPSU_INSTALL_MAGICK=1` or `TPSU_MAGICK=/path/to/magick`.

## Image conversion (UIC)

Run **`uic`**, then **select images** or **select folder**. Outputs stay beside the originals (or in the chosen folder for mass convert).

`tpsu` detects each file’s type from its extension, asks which **format to convert to**, then whether to **delete originals**, an optional **basename** when converting one file, and **overwrite** when the destination already exists.

Conversions use a **PNG hub** (source → PNG → target). Same-format runs do not re-encode: unchanged if the path is the same, or a byte copy when renaming or moving. Formats: JPEG, PNG, GIF, WebP, SVG (read), HEIC, AVIF, TIFF, BMP, RAW (read). GIF uses the first frame only when converting to another format.

**Mass convert** walks one folder (not recursive), reports success rate, skipped non-images, skipped existing files when you decline overwrite, and up to five random failure names.

HEIC, AVIF, and RAW depend on ImageMagick coders/plugins on your system; unavailable write targets are dimmed in the menu.

Integration contract: [ADDONS.md](https://github.com/Cimmerial/tps/blob/main/ADDONS.md) in the tps repo (`contract_iteration: 1` in `tps-addon.yaml`).

From the **tps-utilities** project, **`u`** (update) pulls this repo and rebuilds `bin/tpsu`, then updates **tps** itself (dev repo `git pull` + rebuild, or GitHub release if you have no tps dev checkout linked).
