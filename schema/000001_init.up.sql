CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE translation_jobs
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    
    -- Исходный файл
    original_file_name VARCHAR(255) NOT NULL,
    original_file_path VARCHAR(500) NOT NULL,
    original_content TEXT,
    source_language VARCHAR(10) NOT NULL,
    
    -- Перевод
    target_language VARCHAR(10) NOT NULL,
    translated_content TEXT,
    translated_file_path VARCHAR(500),
    
    -- Статус и метаданные
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    translation_service VARCHAR(50),
    error_message TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_translation_jobs_user_id ON translation_jobs(user_id);
CREATE INDEX idx_translation_jobs_status ON translation_jobs(status);