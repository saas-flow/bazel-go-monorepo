DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'account_type') THEN
        CREATE TYPE account_types AS ENUM ('USER', 'SERVICE_ACCOUNT');
    END IF;
END $$;
