import argparse
import os
import glob
import random
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
    client = make_grpc_client(LB_ADDR, tls_config=TLS_CONFIG)

    img = cv2.imread(image_path)
    if img is None:
        raise FileNotFoundError(f"[Thread {thread_id}] Image '{image_path}' not found")

    img = cv2.resize(img, (224, 224)).astype(np.float32)
    img = img.transpose((2, 0, 1))[None, ...]  # NHWC → NCHW

    output = client.predict({"0": img}, "resnet50")
    idx = int(np.argmax(output[0]))
    print(f"[Thread {thread_id}] {imagenet_classes[idx]} (from {os.path.basename(image_path)})")

def main():
    parser = argparse.ArgumentParser(
        description="Run ResNet50 inference on random images via OVMS gRPC API, in 10 parallel threads"
    )
    parser.add_argument(
        "image_dir",
        help="Path to the directory containing your .jpg/.jpeg images",
    )
    args = parser.parse_args()

    # 1) Gather all .jpg/.jpeg files in the directory
    patterns = ("*.jpg", "*.jpeg", "*.JPG", "*.JPEG")
    images = []
    for pat in patterns:
        images.extend(glob.glob(os.path.join(args.image_dir, pat)))
    if not images:
        parser.error(f"No JPEG files found in '{args.image_dir}'")

    # 2) Launch 10 parallel workers, each with a random pick
    num_threads = 10
    with ThreadPoolExecutor(max_workers=num_threads) as executor:
        futures = []
        # If you want _unique_ images per thread and you have enough files:
        if len(images) >= num_threads:
            picks = random.sample(images, num_threads)
        else:
            # fallback to allowing duplicates
            picks = [random.choice(images) for _ in range(num_threads)]

        for i, img_path in enumerate(picks):
            futures.append(executor.submit(worker, img_path, i))

        for future in as_completed(futures):
            future.result()  # will raise if any worker errored

if __name__ == "__main__":
    main()

