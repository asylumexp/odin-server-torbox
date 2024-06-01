FROM node:alpine as builder
WORKDIR /build/
RUN npm i -g pnpm
COPY package.json . 
RUN pnpm i
COPY . .
RUN pnpm run build

FROM node:alpine
WORKDIR /app
COPY --from=builder /build/.output /app

CMD ["node", "server/index.mjs"]