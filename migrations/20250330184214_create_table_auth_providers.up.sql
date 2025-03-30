
CREATE TABLE IF NOT EXISTS auth_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE, -- "google", "github", "facebook", "twitter", etc.
    client_id TEXT NOT NULL,
    client_secret TEXT NOT NULL,
    auth_url TEXT NOT NULL,     -- URL untuk request authorization code
    token_url TEXT NOT NULL,    -- URL untuk menukar kode dengan access_token
    user_info_url TEXT NOT NULL, -- URL untuk mendapatkan data user (tambahan)
    scopes TEXT NOT NULL DEFAULT 'openid email profile', -- Scopes yang diperlukan
    redirect_uri TEXT NOT NULL, -- Redirect URI untuk callback
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO auth_providers (name, client_id, client_secret, auth_url, token_url, user_info_url, scopes, redirect_uri)
VALUES
('google', 'GOOGLE_CLIENT_ID', 'GOOGLE_CLIENT_SECRET',
 'https://accounts.google.com/o/oauth2/auth',
 'https://oauth2.googleapis.com/token',
 'https://www.googleapis.com/oauth2/v2/userinfo',
 'openid email profile',
 ''),

('github', 'GITHUB_CLIENT_ID', 'GITHUB_CLIENT_SECRET',
 'https://github.com/login/oauth/authorize',
 'https://github.com/login/oauth/access_token',
 'https://api.github.com/user',
 'user:email',
 ''),

('facebook', 'FACEBOOK_CLIENT_ID', 'FACEBOOK_CLIENT_SECRET',
 'https://www.facebook.com/v14.0/dialog/oauth',
 'https://graph.facebook.com/v14.0/oauth/access_token',
 'https://graph.facebook.com/me?fields=id,name,email,picture',
 'public_profile,email',
 '');
