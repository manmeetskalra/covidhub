FROM node:12

WORKDIR /opt/app

COPY package*.json ./
RUN npm install

ENV MYSQL_HOST=covidhub.cluster-cscg2ncxunhq.us-west-2.rds.amazonaws.com
ENV MYSQL_PORT=3306
ENV MYSQL_USER=admin
ENV MYSQL_PASSWORD=Shopify7224

COPY . .

EXPOSE 8080

CMD [ "npm", "start" ]
