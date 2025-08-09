# Stage 1: Build Vite app
FROM node:20-alpine AS build

WORKDIR /app

# Install dependencies first for better caching
COPY package*.json ./
RUN npm install

# Copy app source
COPY . .

ARG SERVER_PORT
ARG SERVER_URL

ENV VITE_SERVER_PORT=${SERVER_PORT}
ENV VITE_SERVER_URL=${SERVER_URL}

RUN echo "VITE_SERVER_PORT=$VITE_SERVER_PORT" && \
    echo "VITE_SERVER_URL=$VITE_SERVER_URL"


# Build the production files
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:1.27-alpine

# Copy build output to Nginx html directory
COPY --from=build /app/dist /usr/share/nginx/html

# Optional: custom nginx config for SPA routing
COPY --from=build nginx.conf /etc/nginx/conf.d/default.conf

# Expose HTTP port
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]