# container base image
FROM node:16-alpine

WORKDIR /app

# copy and download dependencies
COPY package.json ./
COPY package-lock.json ./
RUN npm ci && npm cache clean --force

# copy sources files into the image
COPY api ./api
COPY tsconfig.json ./tsconfig.json

# build 
RUN npm run build

# expose port 3000
EXPOSE 3000

# start app
CMD ["npm", "start"]