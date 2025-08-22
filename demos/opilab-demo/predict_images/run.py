#!/usr/bin/env python

import argparse
import concurrent.futures
import dataclasses
import glob
import os
import random
import typing

import cv2
import numpy
import ovmsclient

import classes


BASEDIR = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))

LB_ADDR = "172.16.3.200:443"
LB_TLSCERT = f"{BASEDIR}/server.crt"


@dataclasses.dataclass(kw_only=True, frozen=True)
class WorkerResult:
    image_path: str
    thread_id: int
    idx: int

    @property
    def image_basename(self) -> str:
        return os.path.basename(self.image_path)

    @property
    def image_class(self) -> str:
        return classes.imagenet_classes[self.idx]


def find_images(image_dir: str) -> typing.Generator[str, None, None]:
    patterns = ("*.jpg", "*.jpeg", "*.JPG", "*.JPEG")
    for pat in patterns:
        for file in glob.glob(os.path.join(image_dir, pat)):
            yield file


def select_n_images_randomly(images: tuple[str, ...], num: int) -> tuple[str, ...]:
    img = list(images)
    random.shuffle(img)
    if num <= len(img):
        img = img[0:num]
    else:
        img.extend(random.choice(images) for _ in range(num - len(img)))
    return tuple(img)


def worker(image_path: str, thread_id: int) -> WorkerResult:
    """
    Worker function: creates its own client, loads & preprocesses
    the image, runs inference, and prints the predicted class.
    """
    client = ovmsclient.make_grpc_client(
        url=LB_ADDR,
        tls_config={
            "server_cert_path": LB_TLSCERT,
        },
    )

    img = cv2.imread(image_path)
    if img is None:
        raise FileNotFoundError(f"[Thread {thread_id}] Image '{image_path}' not found")

    img = cv2.resize(img, (224, 224)).astype(numpy.float32)
    img = img.transpose((2, 0, 1))[None, ...]  # NHWC → NCHW

    output = client.predict({"0": img}, "resnet50")
    idx = int(numpy.argmax(output[0]))

    return WorkerResult(
        image_path=image_path,
        thread_id=thread_id,
        idx=idx,
    )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Run ResNet50 inference on random images via OVMS gRPC API, in 10 parallel threads"
    )
    parser.add_argument(
        "image_dir",
        default="images",
        help="Path to the directory containing your .jpg/.jpeg images",
        nargs="?",
    )

    args = parser.parse_args()

    args.images = tuple(find_images(args.image_dir))
    if not args.images:
        parser.error(f"No JPEG files found in '{args.image_dir}'")

    return args


def main() -> None:
    args = parse_args()

    num_threads = 20
    with concurrent.futures.ThreadPoolExecutor(max_workers=num_threads) as executor:

        picks = select_n_images_randomly(args.images, num_threads)

        futures = [
            executor.submit(worker, img_path, i) for i, img_path in enumerate(picks)
        ]

        results = [
            future.result() for future in concurrent.futures.as_completed(futures)
        ]

    results.sort(key=lambda o: (o.image_basename, o.thread_id))

    print(
        f"Found {len(args.images)} images in {os.path.abspath(args.image_dir)}. Inference with {len(results)} threads."
    )
    for result in results:
        print(
            f"[Thread {result.thread_id:3}] {result.image_class} (from {result.image_basename})"
        )


if __name__ == "__main__":
    main()
