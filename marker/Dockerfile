FROM python:3.12-slim

WORKDIR /usr/src/app

COPY pyproject.toml ./
COPY marker/ marker/
RUN pip install --no-cache-dir .

# Verify the entry point works
RUN marker --help > /dev/null

CMD ["marker", "--config-file", "/mnt/input.json", "--output-file", "/mnt/output.json"]
