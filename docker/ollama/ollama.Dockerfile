FROM ollama/ollama:0.12.5

RUN ollama serve & \
  sleep 5 && \
  ollama pull all-minilm:22m && \
  pkill ollama

CMD ["serve"]
