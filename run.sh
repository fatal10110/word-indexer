# /bin/bash
docker build -t word-indexer .
docker run --rm -p 8080:8080 word-indexer