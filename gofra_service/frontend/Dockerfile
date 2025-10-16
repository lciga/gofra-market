### Build stage: use node to build the SPA
FROM node:18-alpine AS builder

WORKDIR /app

# install deps
COPY package.json package-lock.json* ./
RUN npm ci --silent --no-audit --no-fund

# Patch: ensure quasar import-transformation shim is available inside installed package
# Some quasar distributions don't include the runtime JS transform file; copy our shim
# into the installed package so the @quasar/app webpack loader can require it.
COPY quasar/dist/transforms/import-transformation.js /app/node_modules/quasar/dist/transforms/import-transformation.js

# copy sources
COPY ./src ./src
COPY quasar.conf.js ./quasar.conf.js
COPY package.json .
COPY . .

# build production bundle
ENV NODE_ENV=production
RUN npx quasar build

### Final stage: serve static with nginx (small runtime)
FROM nginx:stable-alpine

# Remove default nginx config and copy our static build
RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/dist/spa /usr/share/nginx/html

# Minimal nginx config: serve index.html for SPA (history mode)
COPY ./frontend-nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
