# 1. Этап сборки (build stage)
FROM node:18-alpine AS build

# Создадим рабочую директорию
WORKDIR /app

# Скопируем package*.json, установим зависимости
COPY package*.json .
RUN npm install


# Скопируем всё остальное и соберём проект
COPY . .
ENV NODE_ENV=production
ENV VITE_SESSION_TIME_IN_MILLISECONDS=3600000
#ENV VITE_LOGIN_API_URL='http://j-auth-ppak.vmeda.local/api/v2/login'
ENV VITE_API_URL='/api/v1'
RUN npm run build

# 2. Этап запуска (release stage)
FROM nginx:alpine

RUN mkdir static

# Скопируем собственный конфиг Nginx (если у вас есть такой)
# Предположим, что ваш конфиг лежит в корне проекта и называется nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем собранные файлы из первого этапа
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
