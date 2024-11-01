server {
    server_name  gowsay.vnykmshr.com;
    listen 80;

    access_log /var/log/nginx/gowsay.access.log;
    error_log  /var/log/nginx/gowsay.error.log;

    root /var/www/gowsay/public;

    location / {
        proxy_pass http://127.0.0.1:9000;
    }
}
