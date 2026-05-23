FROM node:18-alpine
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install --silent
COPY . .
EXPOSE 8000
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "8000"]
