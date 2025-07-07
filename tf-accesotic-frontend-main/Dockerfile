# spa-client-angular/Dockerfile

# Etapa 1: Compilación de la app Angular
FROM node:18 AS build
WORKDIR /app
COPY . .
RUN npm install -g @angular/cli && npm install && ng build --configuration=production

# Etapa 2: Servir con Nginx
FROM nginx:alpine
COPY --from=build /app/dist/spa-client-angular /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
