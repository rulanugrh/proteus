FROM node:21-alpine

ARG EXPOSE_PORT

# Custom workdir
WORKDIR /app

# Copy all files
COPY . .

# Install depedency
RUN npm i

# Build file TS to binary
RUN npm run build

# EXPOSE App Docker
EXPOSE ${EXPOSE_PORT}

# Running NPM
CMD [ "npm", "start" ]