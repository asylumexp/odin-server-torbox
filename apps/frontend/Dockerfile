FROM node:21 as builder
WORKDIR /build/
COPY . . 
RUN npm i -g pnpm
RUN pnpm i
RUN pnpm run generate


FROM lipanski/docker-static-website:latest
COPY --from=builder /build/.output/. .