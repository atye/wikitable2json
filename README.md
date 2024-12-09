For Go users, you can use the client directly. See the `examples` folder.

# Running Locally
```
docker build -t wikitable2json .
docker run --rm -p 8080:8080 -e PORT=8080 -e CACHE_SIZE=10 -e CACHE_EXPIRATION=10s wikitable2json
```
