FROM python:3.12-slim AS builder

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

FROM python:3.12-slim

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

COPY --from=builder /install /usr/local
COPY app ./app
COPY run.py ./run.py

EXPOSE 8123

HEALTHCHECK --interval=30s --timeout=3s --start-period=20s --retries=5 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8123/api/health/')"

CMD ["python", "run.py", "--env", "prod", "--host", "0.0.0.0", "--port", "8123"]
