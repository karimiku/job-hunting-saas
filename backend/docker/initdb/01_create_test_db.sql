-- テスト用データベースを作成（存在しない場合のみ）
SELECT 'CREATE DATABASE job_hunting_test'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'job_hunting_test')\gexec
