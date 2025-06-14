-- 1. 创建表结构
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(50) NOT NULL,
    status TINYINT DEFAULT 1 COMMENT '1-启用 0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '权限标识符',
    name VARCHAR(50) NOT NULL,
    type ENUM('menu', 'api', 'button') NOT NULL COMMENT '权限类型',
    description VARCHAR(255)
);

-- 关联表
CREATE TABLE IF NOT EXISTS user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);

-- 2. 初始化数据
INSERT INTO permissions (code, name, type, description) VALUES 
('user:list', '用户列表', 'api', '查看用户列表'),
('user:create', '创建用户', 'api', '新增用户权限'),
('role:assign', '分配角色', 'api', '为用户分配角色');

INSERT INTO roles (name, description) VALUES 
('admin', '管理员'),
('editor', '编辑员');

-- 管理员拥有所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- 初始管理员用户（密码：admin123）
INSERT INTO users (username, password, salt, status) VALUES 
('admin', '$2a$10$xVCHq6I7fBzYXLpD7KQZQuT3RlEQj4XgY/6Z7bA8Jk9JKlLd1J1XW', 'abc123', 1);

INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);
