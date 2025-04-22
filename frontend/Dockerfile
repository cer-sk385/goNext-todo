# 開発用
FROM node:18-alpine

WORKDIR /app

# 依存関係のインストール
COPY package*.json ./
RUN npm install

# ソースコードのコピー
COPY . .

# 開発サーバーの起動
EXPOSE 3000
CMD ["npm", "run", "dev"]
