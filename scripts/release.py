#!/usr/bin/env python3
"""Collect verified installer archives under build/release for CI and Releases."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import shutil
import subprocess
import tarfile

ROOT = Path(__file__).resolve().parent.parent
CONFIG = json.loads((ROOT / "mygo.json").read_text())


def collect(target):
    system, arch = target.split("/")
    source = ROOT / CONFIG["out"] / f"{system}-{arch}"
    output = ROOT / CONFIG["out"] / "release"
    output.mkdir(parents=True, exist_ok=True)
    name = f'myclip-{CONFIG["version"]}'
    required = {"LICENSE", "NOTICE.md", "THIRD_PARTY_NOTICES.txt"}
    if system == "linux":
        machine = {"amd64": 62, "arm64": 183}[arch]
        for binary in [source / "myclip", source / "libmygo-sqlite3.so"]:
            header = binary.read_bytes()[:20]
            if header[:6] != b"\x7fELF\x02\x01" or int.from_bytes(header[18:20], "little") != machine:
                raise RuntimeError(f"Wrong ELF architecture for {target}: {binary.name}")
        archive = next(source.glob("*.tar.gz"))
        with tarfile.open(archive) as tf:
            names = {Path(x.name).name for x in tf.getmembers()}
            if not required <= names or "libmygo-sqlite3.so" not in names:
                raise RuntimeError("Portable archive is missing bundled notices or SQLite")
        deb = next(source.glob("*.deb"))
        if subprocess.check_output(["dpkg-deb", "-f", str(deb), "Architecture"], text=True).strip() != arch:
            raise RuntimeError("Debian architecture does not match the build target")
        listing = subprocess.check_output(["dpkg-deb", "--contents", str(deb)], text=True)
        if any(x not in listing for x in required | {"libmygo-sqlite3.so"}):
            raise RuntimeError("Debian package is missing bundled notices or SQLite")
        for path, suffix in [(archive, "tar.gz"), (deb, "deb")]:
            shutil.copy2(path, output / f"{name}-linux-{arch}.{suffix}")
    elif target == "darwin/universal":
        app = next(source.glob("*.app"))
        resources = app / "Contents/Resources"
        if any(not (resources / x).is_file() for x in required | {"libmygo-sqlite3.dylib"}):
            raise RuntimeError("macOS bundle is missing notices or SQLite")
        subprocess.run(["codesign", "--verify", "--deep", "--strict", str(app)], check=True)
        for binary in [app / "Contents/MacOS/myclip", resources / "libmygo-sqlite3.dylib"]:
            arches = subprocess.check_output(["lipo", "-archs", str(binary)], text=True).split()
            if set(arches) != {"arm64", "x86_64"}:
                raise RuntimeError(f"Not universal: {binary.name}: {arches}")
        shutil.copy2(next(source.glob("*.dmg")), output / f"{name}-macos-universal.dmg")
        subprocess.run(["ditto", "-c", "-k", "--sequesterRsrc", "--keepParent", str(app),
                        str(output / f"{name}-macos-universal.zip")], check=True)
    else:
        raise ValueError("Unsupported target: " + target)
    print("Release packages verified:", target)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("target", nargs="?", choices=["linux/amd64", "linux/arm64", "darwin/universal"])
    parser.add_argument("--check-tag")
    parser.add_argument("--checksums", action="store_true")
    args = parser.parse_args()
    if args.check_tag:
        if not re.fullmatch(r"v\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?", args.check_tag) or args.check_tag != "v" + CONFIG["version"]:
            raise SystemExit("Release tag must equal v + mygo.json version")
        print("Release tag matches app version")
    elif args.checksums:
        directory = ROOT / CONFIG["out"] / "release"
        files = sorted(x for x in directory.iterdir() if x.is_file() and x.name != "SHA256SUMS")
        prefix = f'myclip-{CONFIG["version"]}'
        expected = {f"{prefix}-linux-{arch}.{suffix}" for arch in ("amd64", "arm64") for suffix in ("deb", "tar.gz")}
        expected |= {f"{prefix}-macos-universal.{suffix}" for suffix in ("dmg", "zip")}
        if {x.name for x in files} != expected:
            raise SystemExit("Release requires exactly the six expected installers")
        (directory / "SHA256SUMS").write_text("".join(
            hashlib.sha256(x.read_bytes()).hexdigest() + "  " + x.name + "\n" for x in files))
    elif args.target:
        collect(args.target)
    else:
        parser.error("provide a target, --check-tag, or --checksums")
