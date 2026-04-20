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
        image = Image.open(io.BytesIO(image_data))
        model = get_model()
        results = model.predict(source=image, save=False, show=False, verbose=False)

        if not results:
            return {}

        probs = getattr(results[0], "probs", None)
        if probs is None:
            return {}

        top5_raw = getattr(probs, "top5", None)
        top5conf_raw = getattr(probs, "top5conf", None)
        if top5_raw is None or top5conf_raw is None:
            return {}

        type_confidence: Dict[str, float] = {}
        for i in range(5):
            label = int(top5_raw[i])
            fish_type = results[0].names[label]
            confidence = top5conf_raw[i]
            type_confidence[FISH_DICT.get(fish_type)] = float(
                confidence.item() if hasattr(confidence, "item") else confidence
            )

        return type_confidence
    except Exception:
        logger.exception("Error during prediction")
        return {}
