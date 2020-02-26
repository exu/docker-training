FROM node:13.8.0-alpine3.10 as builder
WORKDIR /app
ADD . /app
RUN npm install
RUN npm run build


FROM nginx:1.17.8
COPY --from=builder /app/dist /usr/share/nginx/html
