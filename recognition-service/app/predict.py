import io
import logging
import os
from typing import Dict

from PIL import Image
from ultralytics import YOLO

logger = logging.getLogger(__name__)

MODEL_PATH = os.getenv("MODEL_PATH", "./models/best.pt")

FISH_DICT = {
    "Grass_Carp": "草鱼",
    "Crucian_carp": "鲫鱼",
    "Silver_carp": "鲢鱼",
    "Perch": "鲈鱼",
    "Carp": "鲤鱼",
    "Snakehead": "黑鱼",
    "Herring": "青鱼",
    "Parabramis_pekinensis": "鳊鱼",
    "Eel": "鳗鱼",
    "Tilapia_mossambica": "罗非鱼",
    "Noise": "无鱼",
}

_model = None


def get_model() -> YOLO:
    global _model
    if _model is None:
        logger.info("Loading YOLO model from %s", MODEL_PATH)
        _model = YOLO(MODEL_PATH)
    return _model


def get_model_version() -> str:
    return os.path.basename(MODEL_PATH)


def predict_type_confidence(image_data: bytes) -> Dict[str, float]:
    try:
        image = Image.open(io.BytesIO(image_data)).convert("RGB")
        model = get_model()
        results = model.predict(source=image, save=False, show=False, verbose=False)

        if not results:
            return {}

        probs = getattr(results[0], "probs", None)
        if probs is None:
            return {}

        top5 = list(getattr(probs, "top5", []) or [])[:5]
        top5conf = list(getattr(probs, "top5conf", []) or [])[:5]

        if not top5 or not top5conf:
            return {}

        type_confidence: Dict[str, float] = {}
        for label_idx, confidence in zip(top5, top5conf):
            label_int = int(label_idx)
            english_name = results[0].names.get(label_int, str(label_int))
            fish_type = FISH_DICT.get(english_name, english_name)
            conf_value = float(confidence.item() if hasattr(confidence, "item") else confidence)
            type_confidence[fish_type] = conf_value

        return type_confidence
    except Exception:
        logger.exception("Error during prediction")
        return {}
