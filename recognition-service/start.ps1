$ErrorActionPreference = 'Stop'

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptDir

if (!(Test-Path .venv)) {
  python -m venv .venv
}

. .\.venv\Scripts\Activate.ps1
pip install -r requirements.txt

$host = if ($env:HOST) { $env:HOST } else { '0.0.0.0' }
$port = if ($env:PORT) { $env:PORT } else { '8001' }

uvicorn app.main:app --host $host --port $port
