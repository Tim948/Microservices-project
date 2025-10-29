-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы проектов
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    owner_id INTEGER REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы задач
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    project_id INTEGER REFERENCES projects(id),
    assigned_to INTEGER REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'pending',
    priority VARCHAR(20) DEFAULT 'medium',
    due_date TIMESTAMP,
    estimated_hours DECIMAL(5,2),
    actual_hours DECIMAL(5,2),
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы уведомлений
CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    related_entity_type VARCHAR(50),
    related_entity_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы активности пользователей
CREATE TABLE IF NOT EXISTS user_activities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    activity_type VARCHAR(100) NOT NULL,
    description TEXT,
    entity_type VARCHAR(50),
    entity_id INTEGER,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_activities_user_id ON user_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activities_created_at ON user_activities(created_at);

-- Вставка тестовых данных
INSERT INTO users (username, email, first_name, last_name, role) VALUES
('admin', 'admin@company.com', 'Admin', 'User', 'admin'),
('john_doe', 'john@company.com', 'John', 'Doe', 'manager'),
('jane_smith', 'jane@company.com', 'Jane', 'Smith', 'user'),
('mike_wilson', 'mike@company.com', 'Mike', 'Wilson', 'user');

INSERT INTO projects (name, description, owner_id) VALUES
('Website Redesign', 'Complete redesign of company website', 1),
('Mobile App Development', 'Development of new mobile application', 2),
('Marketing Campaign', 'Q4 marketing campaign preparation', 2);

INSERT INTO tasks (title, description, project_id, assigned_to, status, priority, due_date, estimated_hours, created_by) VALUES
('Design Homepage', 'Create new homepage design', 1, 3, 'in_progress', 'high', '2024-02-15', 8.0, 1),
('Develop API', 'Build REST API for mobile app', 2, 4, 'pending', 'high', '2024-02-20', 16.0, 2),
('Write Content', 'Create marketing content for campaign', 3, 3, 'completed', 'medium', '2024-02-10', 4.0, 2),
('User Testing', 'Conduct user testing for new features', 2, 4, 'pending', 'medium', '2024-02-25', 6.0, 2);