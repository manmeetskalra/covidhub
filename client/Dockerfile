FROM node:13.12.0-alpine

# set working directory
WORKDIR /app

COPY package*.json ./
RUN npm install

ENV REACT_APP_API_URI="http://ec2co-ecsel-1jni8txjq1cve-122992530.us-west-2.elb.amazonaws.com:8080"
ENV REACT_APP_EMAIL="http://ec2co-ecsel-1jni8txjq1cve-122992530.us-west-2.elb.amazonaws.com:8001"
ENV REACT_APP_TEXT="http://ec2co-ecsel-1jni8txjq1cve-122992530.us-west-2.elb.amazonaws.com:8002"

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
