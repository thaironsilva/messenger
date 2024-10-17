-- migration up for create_messages_table
CREATE TABLE messages (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    sender_id uuid NOT NULL,
    receiver_id uuid NOT NULL,
    body VARCHAR(255) NOT NULL CHECK (body <> ''),
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_messages_sender FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_receiver FOREIGN KEY(receiver_id) REFERENCES users(id) ON DELETE CASCADE
);