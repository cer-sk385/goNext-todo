#!/bin/bash

# 色の設定
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}コンテナを停止して削除します...${NC}"
podman-compose down

echo -e "${YELLOW}PostgreSQLのボリュームを削除します...${NC}"
podman volume rm aks-test_postgres_data || true

echo -e "${YELLOW}db ディレクトリが存在しない場合は作成します...${NC}"
mkdir -p db

echo -e "${YELLOW}todos テーブルの初期化スクリプトが存在するか確認します...${NC}"
if [ ! -f db/init.sql ]; then
  echo -e "${YELLOW}init.sql を作成します...${NC}"
  cat > db/init.sql << 'EOF'
CREATE TABLE IF NOT EXISTS todos (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  completed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 初期データ
INSERT INTO todos (title, description) VALUES 
  ('最初のタスク', 'これは最初のTODOタスクです'),
  ('メールを確認', '重要なメールを確認してください');
EOF
  echo -e "${GREEN}init.sql を作成しました${NC}"
fi

echo -e "${YELLOW}コンテナをバックグラウンドで再起動します...${NC}"
podman-compose up -d --build

echo -e "${GREEN}完了しました！${NC}"
echo -e "${GREEN}フロントエンド: http://localhost:3000${NC}"
echo -e "${GREEN}バックエンドAPI: http://localhost:8080${NC}" 