# OPI test frmework

will be using pytest

## to run the test cases

```Shell
docker run --rm --network host \
    --mount src=./,target=/opi/test-cases,type=bind \
    --it opi/test-framework bash
pytest /opi/test-cases/test_sample.py
```
