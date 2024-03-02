FROM node:alpine as builder
WORKDIR /build/
COPY . . 
RUN npm i -g pnpm
RUN pnpm i
RUN pnpm run build

FROM node:alpine
WORKDIR /app
COPY --from=builder /build/.output /app

CMD ["node", "server/index.mjs"]