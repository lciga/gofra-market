FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci --silent --no-audit --no-fund

COPY quasar/dist/transforms/import-transformation.js /app/node_modules/quasar/dist/transforms/import-transformation.js

COPY ./src ./src
COPY quasar.conf.js ./quasar.conf.js
COPY package.json .
COPY . .

ENV NODE_ENV=production
RUN npx quasar build

FROM nginx:stable-alpine

RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/dist/spa /usr/share/nginx/html

COPY ./frontend-nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
