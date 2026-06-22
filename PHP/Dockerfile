FROM php:8.3-apache

RUN a2enmod rewrite

RUN echo "ServerName localhost" > /etc/apache2/conf-available/servername.conf \
    && a2enconf servername


COPY . /var/www/html

WORKDIR /var/www/html
