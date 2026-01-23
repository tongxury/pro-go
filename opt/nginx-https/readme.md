1. 修改domain 删除 [certbot](certbot)目录
2. 配置dns 到nginx 记得加上
   server {
   listen 80;
   server_name resource.veogoai.com;

   # Certbot 验证目录
   location /.well-known/acme-challenge/ {
   root /var/www/certbot;
   }

   # 将所有 HTTP 请求重定向到 HTTPS
   location / {
   return 301 https://$host$request_uri;
   }
   }

2. docker compose up -d nginx 只启动nginx

3. docker run -it --rm --name certbot \
-v "/root/nginx/certbot/etc/letsencrypt:/etc/letsencrypt" \
-v "/root/nginx/certbot/var/lib/letsencrypt:/var/lib/letsencrypt" \
-v "/root/nginx/certbot/var/www/certbot:/var/www/certbot" \
   docker.m.daocloud.io/certbot/certbot certonly --webroot --webroot-path  /var/www/certbot \
--email tongxurt@gmail.com \
-d veogo.cn \
--agree-tos \
--no-eff-email

成功之后是这样
Successfully received certificate.
Certificate is saved at: /etc/letsencrypt/live/i-beta.yoozyai.com/fullchain.pem
Key is saved at:         /etc/letsencrypt/live/i-beta.yoozyai.com/privkey.pem
This certificate expires on 2026-03-16.
These files will be updated when the certificate renews.

NEXT STEPS:
- The certificate will need to be renewed before it expires. Certbot can automatically renew the certificate in the background, but you may need to take steps to enable that functionality. See https://certbot.org/renewal-setup for instructions.

- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
If you like Certbot, please consider supporting our work by:
* Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
* Donating to EFF:                    https://eff.org/donate-le
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

4. 加上 ssl 配置

server {
listen 443 ssl;
server_name asset.veogo-src.com;  # 替换为你的域名

    # SSL 配置
    ssl_certificate /etc/letsencrypt/live/resource.veogoai.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/resource.veogoai.com/privkey.pem;

    # SSL 优化
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;

    # 应用程序配置
location / {
proxy_set_header Host            $host;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_buffering off;
proxy_pass http://adb1f4ba7b808404099bc618f05cbb57-2072134424.ap-southeast-1.elb.amazonaws.com;
}
}