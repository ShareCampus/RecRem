curl https://api.openai.com/v1/embeddings \
  -H "Authorization: Bearer <your key>" \
  -H "Content-Type: application/json" \
  -d '{
    "input": "The food was delicious and the waiter...",
    "model": "text-embedding-3-small",
    "encoding_format": "float"
  }'
