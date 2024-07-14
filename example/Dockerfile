FROM nginx

COPY index.html /usr/share/nginx/html

ARG backend

RUN sed -i "s/backend/${backend}/g" /usr/share/nginx/html/index.html
