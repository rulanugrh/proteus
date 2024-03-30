# FROM to use image, i use golang:1.21.6
# use base alpine because it is lightweight image
FROM golang:1.21.6-alpine

# arg for depend port app into env
ARG HTTP_PORT

# ustom default workdir for running app
WORKDIR /app

# copy all code from local into docker image
COPY . .

# running install modules
RUN go mod install

# running build main file
RUN go build -o build/api

# EXPOSE port to outside Docker Container
EXPOSE ${HTTP_PORT}

# CMD running for binary api
CMD [ "./build/api" ]