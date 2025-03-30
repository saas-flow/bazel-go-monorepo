
CREATE TABLE password_rules (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      version_id UUID NOT NULL,
      name TEXT NOT NULL,
      display_name TEXT NOT NULL,
      pattern TEXT NOT NULL,
      FOREIGN KEY (version_id) REFERENCES password_versions(id) ON DELETE CASCADE
  );