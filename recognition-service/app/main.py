import logging
from typing import Any, Dict

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from app.predict import get_model, get_model_version, predict_type_confidence

logging.basicConfig(level=logging.INFO)

app = FastAPI(title="Smart Fish Recognition Service", version="1.0.0")


@app.get("/healthz")
def healthz() -> Dict[str, Any]:
    try:
        get_model()
        return {"status": "ok", "model_loaded": True, "model_version": get_model_version()}
    except Exception as e:
        return {"status": "degraded", "model_loaded": False, "error": str(e)}


@app.post("/predict/fish-species")
async def predict_fish_species(request: Request):
    form = await request.form()
    file = form.get("fish_pic")
    if file is None:
        return JSONResponse({"error": "missing fish_pic"}, status_code=400)

    image_data = await file.read()
    result = predict_type_confidence(image_data)

    if not result:
        return {"type_confidence": {}, "model_version": get_model_version()}

    return {"type_confidence": result, "model_version": get_model_version()}


@app.post("/yolo")
async def yolo_compat(request: Request):
    form = await request.form()
    if "fish_pic" not in form:
        return JSONResponse({"msg": "No file part in the request"}, status_code=400)

    file = form["fish_pic"]
    if getattr(file, "filename", "") == "":
        return JSONResponse({"msg": "No file selected for uploading"}, status_code=400)

    image_data = await file.read()
    result = predict_type_confidence(image_data)

    if result:
        fish_type = max(result, key=result.get)
        fish_confidence = result[fish_type]
        return {
            "msg": "Classification succeed",
            "type": fish_type,
            "confidence": fish_confidence,
        }

    return {"msg": "Failed to classify"}
