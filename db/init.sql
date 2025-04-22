CREATE TABLE IF NOT EXISTS todos (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  completed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 初期データ（必要な場合）
INSERT INTO todos (title, description) VALUES 
  ('最初のタスク', 'これは最初のTODOタスクです'),
  ('メールを確認', '重要なメールを確認してください');
