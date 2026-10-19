FROM node:18-alpine

WORKDIR /app

COPY . .

# node_modules di-install di host dan di-bind-mount lewat docker-compose,
# jadi tidak perlu di-install di dalam image.

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host"]
