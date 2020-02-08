FROM php:7-apache

RUN pecl install -o -f redis \
    &&  rm -rf /tmp/pear \
    &&  docker-php-ext-enable redis

ADD index.php /var/www/html/
