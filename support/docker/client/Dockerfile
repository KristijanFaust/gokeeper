FROM node:16.3.0-alpine

WORKDIR /gokeeper-ui

COPY ui/public ./public
COPY ui/src ./src
COPY ui/package.json ./
COPY ui/package-lock.json ./

RUN npm install

EXPOSE 3000

CMD [ "npm", "start" ]