# Recognition Service

YOLO/Ultralytics fish recognition service for Smart Fish.

## Features

- `POST /predict/fish-species`: returns `type_confidence` map for Go backend
- `POST /yolo`: Flask-compatible response (`msg`, `type`, `confidence`)
- `GET /healthz`: service and model status

## Requirements

- Python 3.10+
- A trained Ultralytics model file at `./models/best.pt`

## Setup

```bash
cd recognition-service
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

Windows PowerShell:

```powershell
cd recognition-service
python -m venv .venv
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
```

## Run

```bash
cd recognition-service
uvicorn app.main:app --host 0.0.0.0 --port 8001
```

## Environment Variables

- `MODEL_PATH` (default: `./models/best.pt`)
- `HOST` (default: `0.0.0.0`)
- `PORT` (default: `8001`)

## Go Backend Integration

Go backend endpoint: `POST /api/v1/yolo` (Flask-compat).

Optional backend env vars:

- `YOLO_INFER_URL` (default: `http://127.0.0.1:8001/predict/fish-species`)
- `YOLO_INFER_TIMEOUT_SEC` (default: `5`)

If you set `YOLO_INFER_URL` to `/yolo`, backend is still compatible.
