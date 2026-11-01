package native

import (
	"encoding/base64"
	"unicode/utf16"
)

// The Windows OCR engine is a WinRT API. Rather than wire WinRT interop
// into Go — the way Linux shells out to tesseract, Windows shells out to
// PowerShell, which speaks WinRT natively. The script reads the image
// path from the environment and prints the recognized text on stdout.
const ocrScript = `$ErrorActionPreference = 'Stop'
$null = [Windows.Graphics.Imaging.BitmapDecoder, Windows.Foundation, ContentType = WindowsRuntime]
$null = [Windows.Media.Ocr.OcrEngine, Windows.Foundation, ContentType = WindowsRuntime]
$null = [Windows.Storage.Streams.InMemoryRandomAccessStream, Windows.Foundation, ContentType = WindowsRuntime]
$asTask = ([System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object {
	$_.Name -eq 'AsTask' -and $_.GetParameters().Count -eq 1 -and
	$_.GetParameters()[0].ParameterType.Name -eq 'IAsyncOperation' + [char]0x60 + '1' })[0]
function Await($op, [type]$t) {
	$task = $asTask.MakeGenericMethod($t).Invoke($null, @($op))
	$task.Wait(-1) | Out-Null
	$task.Result
}
$bytes = [System.IO.File]::ReadAllBytes($env:MYCLIP_OCR_INPUT)
$stream = New-Object Windows.Storage.Streams.InMemoryRandomAccessStream
$writer = New-Object Windows.Storage.Streams.DataWriter($stream.GetOutputStreamAt(0))
$writer.WriteBytes($bytes)
$writer.StoreAsync() | Out-Null
$writer.FlushAsync() | Out-Null
$stream.Seek(0) | Out-Null
$decoder = Await ([Windows.Graphics.Imaging.BitmapDecoder]::CreateAsync($stream)) ([Windows.Graphics.Imaging.BitmapDecoder])
$bitmap = Await ($decoder.GetSoftwareBitmapAsync()) ([Windows.Graphics.Imaging.SoftwareBitmap])
$engine = [Windows.Media.Ocr.OcrEngine]::TryCreateFromUserProfileLanguages()
if (-not $engine) {
	$engine = [Windows.Media.Ocr.OcrEngine]::TryCreateFromLanguage((New-Object Windows.Globalization.Language('en-US')))
}
if (-not $engine) { exit 3 }
$result = Await ($engine.RecognizeAsync($bitmap)) ([Windows.Media.Ocr.OcrResult])
[Console]::Out.Write($result.Text)
`

// powershellCommand encodes the script the way -EncodedCommand wants:
// base64 of UTF-16LE. Pure, so the encoding is tested everywhere.
func powershellCommand() string {
	u := utf16.Encode([]rune(ocrScript))
	b := make([]byte, 2*len(u))
	for i, v := range u {
		b[i*2], b[i*2+1] = byte(v), byte(v>>8)
	}
	return base64.StdEncoding.EncodeToString(b)
}
