-- URL Shortener Service Database Schema
-- PostgreSQL Database Setup Script

-- Create database (run this separately or through createdb command)
-- CREATE DATABASE url_shortener;

-- Connect to the database
-- \c url_shortener;

-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    plan_type VARCHAR(20) DEFAULT 'free' CHECK (plan_type IN ('free', 'pro', 'enterprise')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Create custom domains table
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    domain_name VARCHAR(255) UNIQUE NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    verification_token VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP WITH TIME ZONE
);

-- Create urls table (main table for shortened URLs)
CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    original_url TEXT NOT NULL,
    short_code VARCHAR(20) UNIQUE NOT NULL,
    custom_alias VARCHAR(50),
    domain_id UUID REFERENCES domains(id) ON DELETE SET NULL,
    title VARCHAR(500),
    description TEXT,
    password_hash VARCHAR(255), -- for password-protected URLs
    is_active BOOLEAN DEFAULT true,
    click_count INTEGER DEFAULT 0,
    unique_click_count INTEGER DEFAULT 0,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_clicked_at TIMESTAMP WITH TIME ZONE
);

-- Create clicks table for detailed analytics
CREATE TABLE clicks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url_id UUID NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    country VARCHAR(2), -- ISO country code
    region VARCHAR(100),
    city VARCHAR(100),
    browser VARCHAR(50),
    os VARCHAR(50),
    device_type VARCHAR(20) CHECK (device_type IN ('desktop', 'mobile', 'tablet', 'unknown')),
    is_unique BOOLEAN DEFAULT false,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create qr_codes table for storing QR code information
CREATE TABLE qr_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url_id UUID NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    qr_data TEXT NOT NULL, -- base64 encoded QR code image
    format VARCHAR(10) DEFAULT 'png' CHECK (format IN ('png', 'jpg', 'svg')),
    size INTEGER DEFAULT 200,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create rate_limits table for tracking API usage
CREATE TABLE rate_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    api_key VARCHAR(64),
    ip_address INET,
    endpoint VARCHAR(200),
    request_count INTEGER DEFAULT 1,
    window_start TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance

-- Users table indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_api_key ON users(api_key);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Domains table indexes
CREATE INDEX idx_domains_user_id ON domains(user_id);
CREATE INDEX idx_domains_domain_name ON domains(domain_name);
CREATE INDEX idx_domains_is_active ON domains(is_active);

-- URLs table indexes
CREATE INDEX idx_urls_user_id ON urls(user_id);
CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_custom_alias ON urls(custom_alias);
CREATE INDEX idx_urls_is_active ON urls(is_active);
CREATE INDEX idx_urls_expires_at ON urls(expires_at);
CREATE INDEX idx_urls_created_at ON urls(created_at);
CREATE INDEX idx_urls_click_count ON urls(click_count);

-- Clicks table indexes
CREATE INDEX idx_clicks_url_id ON clicks(url_id);
CREATE INDEX idx_clicks_clicked_at ON clicks(clicked_at);
CREATE INDEX idx_clicks_country ON clicks(country);
CREATE INDEX idx_clicks_ip_address ON clicks(ip_address);
CREATE INDEX idx_clicks_device_type ON clicks(device_type);
CREATE INDEX idx_clicks_is_unique ON clicks(is_unique);

-- QR codes table indexes
CREATE INDEX idx_qr_codes_url_id ON qr_codes(url_id);

-- Rate limits table indexes
CREATE INDEX idx_rate_limits_user_id ON rate_limits(user_id);
CREATE INDEX idx_rate_limits_api_key ON rate_limits(api_key);
CREATE INDEX idx_rate_limits_ip_address ON rate_limits(ip_address);
CREATE INDEX idx_rate_limits_window_start ON rate_limits(window_start);

-- Create composite indexes for common queries
CREATE INDEX idx_urls_user_active ON urls(user_id, is_active);
CREATE INDEX idx_clicks_url_date ON clicks(url_id, clicked_at);
CREATE INDEX idx_clicks_unique_url ON clicks(url_id, is_unique);

-- Create functions for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for automatic timestamp updates
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_urls_updated_at
    BEFORE UPDATE ON urls
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to generate short codes
CREATE OR REPLACE FUNCTION generate_short_code(length INTEGER DEFAULT 6)
RETURNS TEXT AS $$
DECLARE
    chars TEXT := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result TEXT := '';
    i INTEGER;
BEGIN
    FOR i IN 1..length LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::INTEGER, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Create views for analytics

-- Daily clicks view
CREATE OR REPLACE VIEW daily_clicks AS
SELECT
    DATE(clicked_at) as click_date,
    url_id,
    COUNT(*) as total_clicks,
    COUNT(CASE WHEN is_unique THEN 1 END) as unique_clicks,
    COUNT(DISTINCT country) as countries_count,
    COUNT(DISTINCT device_type) as device_types_count
FROM clicks
GROUP BY DATE(clicked_at), url_id;

-- User statistics view
CREATE OR REPLACE VIEW user_stats AS
SELECT
    u.id as user_id,
    u.email,
    u.plan_type,
    COUNT(urls.id) as total_urls,
    COUNT(CASE WHEN urls.is_active THEN 1 END) as active_urls,
    COALESCE(SUM(urls.click_count), 0) as total_clicks,
    COALESCE(SUM(urls.unique_click_count), 0) as total_unique_clicks,
    u.created_at as user_created_at
FROM users u
LEFT JOIN urls ON u.id = urls.user_id
GROUP BY u.id, u.email, u.plan_type, u.created_at;

-- Top URLs view
CREATE OR REPLACE VIEW top_urls AS
SELECT
    u.id,
    u.short_code,
    u.original_url,
    u.title,
    u.click_count,
    u.unique_click_count,
    u.created_at,
    us.email as user_email
FROM urls u
LEFT JOIN users us ON u.user_id = us.id
WHERE u.is_active = true
ORDER BY u.click_count DESC;

-- Insert sample data for development
INSERT INTO users (email, password_hash, api_key, first_name, last_name, plan_type) VALUES
('john.doe@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeVMpYlqhRVtCSpeW', 'api_key_john_123456789', 'John', 'Doe', 'free'),
('jane.smith@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeVMpYlqhRVtCSpeW', 'api_key_jane_987654321', 'Jane', 'Smith', 'pro'),
('admin@urlshortener.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LeVMpYlqhRVtCSpeW', 'api_key_admin_admin123', 'Admin', 'User', 'enterprise');

-- Insert sample URLs
INSERT INTO urls (user_id, original_url, short_code, title, custom_alias) VALUES
((SELECT id FROM users WHERE email = 'john.doe@example.com'), 'https://www.google.com', 'abc123', 'Google Search', 'google'),
((SELECT id FROM users WHERE email = 'jane.smith@example.com'), 'https://www.github.com', 'xyz789', 'GitHub', 'github'),
((SELECT id FROM users WHERE email = 'john.doe@example.com'), 'https://www.stackoverflow.com', 'def456', 'Stack Overflow', null);

-- Add some sample clicks
INSERT INTO clicks (url_id, ip_address, user_agent, country, region, city, browser, os, device_type, is_unique) VALUES
((SELECT id FROM urls WHERE short_code = 'abc123'), '192.168.1.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', 'US', 'California', 'San Francisco', 'Chrome', 'Windows', 'desktop', true),
((SELECT id FROM urls WHERE short_code = 'abc123'), '192.168.1.2', 'Mozilla/5.0 (iPhone; CPU iPhone OS)', 'US', 'New York', 'New York', 'Safari', 'iOS', 'mobile', true),
((SELECT id FROM urls WHERE short_code = 'xyz789'), '10.0.0.1', 'Mozilla/5.0 (Macintosh; Intel Mac OS X)', 'CA', 'Ontario', 'Toronto', 'Firefox', 'macOS', 'desktop', true);

-- Update click counts in URLs table (this would normally be done by triggers or application logic)
UPDATE urls SET
    click_count = (SELECT COUNT(*) FROM clicks WHERE clicks.url_id = urls.id),
    unique_click_count = (SELECT COUNT(*) FROM clicks WHERE clicks.url_id = urls.id AND is_unique = true);

-- Create stored procedures for common operations

-- Procedure to increment click count
CREATE OR REPLACE FUNCTION increment_click_count(url_short_code TEXT)
RETURNS VOID AS $$
BEGIN
    UPDATE urls
    SET
        click_count = click_count + 1,
        last_clicked_at = CURRENT_TIMESTAMP
    WHERE short_code = url_short_code AND is_active = true;
END;
$$ LANGUAGE plpgsql;

-- Procedure to clean expired URLs
CREATE OR REPLACE FUNCTION cleanup_expired_urls()
RETURNS INTEGER AS $$
DECLARE
    expired_count INTEGER;
BEGIN
    UPDATE urls
    SET is_active = false
    WHERE expires_at < CURRENT_TIMESTAMP AND is_active = true;

    GET DIAGNOSTICS expired_count = ROW_COUNT;
    RETURN expired_count;
END;
$$ LANGUAGE plpgsql;

-- Display table information
SELECT 'Database schema created successfully!' as status;
