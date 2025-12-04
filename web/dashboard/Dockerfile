# Stage 1: Build the Angular application
FROM node:22-alpine AS build

WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install

COPY . .
RUN pnpm run build

# Stage 2: Serve the application with Nginx
FROM nginx:alpine

COPY --from=build /app/dist/dashboard/browser /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
