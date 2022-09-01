FROM node:lts-alpine

WORKDIR /front

COPY package.json .
COPY yarn.lock .

## RUN yarn

EXPOSE 7000

CMD ["yarn", "dev"]