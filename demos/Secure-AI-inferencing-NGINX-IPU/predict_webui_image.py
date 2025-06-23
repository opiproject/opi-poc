import cv2
import numpy as np
import gradio as gr
from ovmsclient import make_grpc_client
from classes import imagenet_classes

# point this at your OVMS load-balancer
LB_ADDR = "1.2.3.4:443"
TLS_CONFIG = {"server_cert_path": "/root/summit/server.crt"}

# init once
client = make_grpc_client(LB_ADDR, tls_config=TLS_CONFIG)

def predict(img: np.ndarray):
    # img comes in as H×W×3 uint8 from Gradio
    img = cv2.resize(img, (224, 224)).astype(np.float32)
    img = img.transpose((2, 0, 1))[None, ...]
    output = client.predict({"0": img}, "resnet50")
    idx = int(np.argmax(output[0]))
    return imagenet_classes[idx]

iface = gr.Interface(
    fn=predict,
    inputs=gr.Image(type="numpy"),
    outputs="text",
    title="ResNet50 on OVMS",
    description="Upload an image and get back its label from ResNet50."
)

if __name__ == "__main__":
    iface.launch(server_name="127.0.0.1", server_port=7860, share=True)

