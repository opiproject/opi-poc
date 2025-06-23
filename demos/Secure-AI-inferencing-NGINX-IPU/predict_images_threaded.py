import argparse
import numpy as np
import cv2
from classes import imagenet_classes
from ovmsclient import make_grpc_client
from concurrent.futures import ThreadPoolExecutor, as_completed

LB_ADDR = "172.16.3.200:443"
TLS_CONFIG = {
    "server_cert_path": "/root/summit/server.crt",
}

def worker(image_path: str, thread_id: int):
    """
    Worker function: creates its own client, loads & preprocesses
    the image, runs inference, and prints the predicted class.
    """
    # 1. Create the gRPC client
    client = make_grpc_client(LB_ADDR, tls_config=TLS_CONFIG)

    # 2. Load & preprocess the image
    img = cv2.imread(image_path)
    if img is None:
        raise FileNotFoundError(f"[Thread {thread_id}] Image '{image_path}' not found or could not be opened")

    img = cv2.resize(img, (224, 224)).astype(np.float32)
    img = img.transpose((2, 0, 1))[None, ...]  # NHWC → NCHW

    # 3. Run inference
    output = client.predict({"0": img}, "resnet50")
    idx = int(np.argmax(output[0]))
    print(f"[Thread {thread_id}] {imagenet_classes[idx]}")

def main():
    parser = argparse.ArgumentParser(
        description="Run ResNet50 inference on an image via OVMS gRPC API, in 10 parallel threads"
    )
    parser.add_argument(
        "image_path",
        help="Path to the input image file",
    )
    args = parser.parse_args()

    # Launch 10 parallel workers
    num_threads = 10
    with ThreadPoolExecutor(max_workers=num_threads) as executor:
        futures = [
            executor.submit(worker, args.image_path, i)
            for i in range(num_threads)
        ]
        # Wait for all to complete (and re-raise any exceptions)
        for future in as_completed(futures):
            future.result()

if __name__ == "__main__":
    main()

