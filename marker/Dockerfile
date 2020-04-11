FROM python:3

WORKDIR /usr/src/app
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
COPY scripts ./scripts/

ENV PATH="/usr/local/bin/python:${PATH}"

CMD [ "/usr/local/bin/python", "-m", "scripts", "--config-file", "/mnt/input.json" ]
