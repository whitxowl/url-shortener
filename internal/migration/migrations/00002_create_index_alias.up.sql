-- +goose Up
CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);