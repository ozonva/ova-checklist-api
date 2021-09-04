-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS checklists (
    user_id         BIGINT NOT NULL,
    checklist_id    UUID NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    data            JSONB,
    PRIMARY KEY (user_id, checklist_id)
);

CREATE INDEX IF NOT EXISTS checklist_id_idx ON checklists (checklist_id);
CREATE INDEX IF NOT EXISTS created_at_idx ON checklists (checklist_id);
-- +goose StatementEnd
