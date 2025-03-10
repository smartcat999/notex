-- Create drafts table
CREATE TABLE IF NOT EXISTS drafts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    summary TEXT,
    category_id INTEGER REFERENCES categories(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on user_id for faster lookups
CREATE INDEX drafts_user_id_idx ON drafts(user_id);

-- Create index on deleted_at for soft deletes
CREATE INDEX drafts_deleted_at_idx ON drafts(deleted_at);

-- Create draft_tags table for many-to-many relationship
CREATE TABLE IF NOT EXISTS draft_tags (
    draft_id INTEGER REFERENCES drafts(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (draft_id, tag_id)
);

-- Create indexes for draft_tags
CREATE INDEX draft_tags_draft_id_idx ON draft_tags(draft_id);
CREATE INDEX draft_tags_tag_id_idx ON draft_tags(tag_id); 