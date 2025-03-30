
ALTER TABLE accounts
DROP COLUMN IF EXISTS auth_provider_id,
DROP COLUMN IF EXISTS provider_user_id,
DROP COLUMN IF EXISTS email_verified