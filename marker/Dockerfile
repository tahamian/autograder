FROM python:3


ADD scripts /usr/local/
COPY requirements.txt /home
RUN echo 'alias marker=python -m /usr/local/scripts' >> ~/.bashrc

ENTRYPOINT []