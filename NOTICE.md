# NOTICE

MyClip is a Go implementation of a clipboard manager for macOS whose feature
set targets parity with [Maccy](https://github.com/p0deje/Maccy) by Alex
Rodionov and contributors, released under the MIT License. MyClip consults
Maccy's documentation, source, and behavior as its specification; MyClip does
not distribute Maccy source code. All credit for the original design belongs
to the Maccy authors.

Maccy's MIT license notice:

> The MIT License (MIT)
>
> Copyright (c) 2025 Alex Rodionov
>
> Permission is hereby granted, free of charge, to any person obtaining a copy
> of this software and associated documentation files (the "Software"), to deal
> in the Software without restriction, including without limitation the rights
> to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
> copies of the Software, and to permit persons to whom the Software is
> furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included in
> all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
> IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
> FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
> AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
> LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
> OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
> THE SOFTWARE.

Lucide icons are used under the ISC license, with the upstream Feather MIT
notice where applicable; see licenses/Lucide.txt. The initial application
icon comes from MyGo's MIT-licensed native template.

Verbatim third-party notices for the pinned runtime dependencies (MyGo,
purego, go-text, x/image, x/text, Go and SQLite) are collected in
THIRD_PARTY_NOTICES.txt. Regenerate them with `python3 scripts/licenses.py`;
CI rejects stale notices. All installers include LICENSE, this notice,
and THIRD_PARTY_NOTICES.txt.
