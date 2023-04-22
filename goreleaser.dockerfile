FROM node:19-alpine3.16 AS remix
WORKDIR /src/github.com/frantjc/sneasler
COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile
COPY . .
RUN yarn build

FROM node:19-alpine3.16
ENTRYPOINT ["sneasler"]
ENV SNEASLER_JS_ENTRYPOINT /app/index.js
COPY sneasler /usr/local/bin
COPY --from=remix /src/github.com/frantjc/sneasler/dist /app
