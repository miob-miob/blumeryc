FROM node:18-alpine3.15 as build

WORKDIR /app

# first copy only dependency files for caching libraries
COPY package.json /app/package.json
COPY package-lock.json /app/package-lock.json

RUN npm install

# COPY . . 
COPY src src
COPY build-scripts build-scripts
COPY webpack.config.ts webpack.config.ts
COPY .env .env
COPY tsconfig.json tsconfig.json

RUN npm run build



# ---------------------------
FROM node:18-alpine3.15 

WORKDIR /app

COPY --from=build /app/dist/bundle /app

RUN rm -rf /app/node_modules

ENTRYPOINT [ "node" ] 
CMD [ "/app/index.js" ]
