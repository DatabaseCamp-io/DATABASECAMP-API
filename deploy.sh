TAG=$1
docker build -t ganinw13120/dbc-backend:$TAG .
docker tag ganinw13120/dbc-backend:$TAG ganinw13120/dbc-backend:$TAG
docker push ganinw13120/dbc-backend:$TAG

# docker pull ganinw13120/dbc-backend:$1
# docker stop backend
# docker run --rm --name backend -p 8080:8080 -d ganinw13120/dbc-backend:$1