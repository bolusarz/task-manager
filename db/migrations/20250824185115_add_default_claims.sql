-- +goose Up
-- +goose StatementBegin
INSERT INTO claims (name) VALUES 
('can_create_task'),
('can_edit_task'),
('can_delete_task'),
('can_assign_task'),
('can_create_team'),
('can_edit_team'),
('can_delete_team'),
('can_manage_team_members'),
('can_manage_roles'),
('can_manage_claims');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM claims WHERE name IN (
    'can_create_task',
    'can_edit_task',
    'can_delete_task',
    'can_assign_task',
    'can_create_team',
    'can_edit_team',
    'can_delete_team',
    'can_manage_team_members',
    'can_manage_roles',
    'can_manage_claims'
);
-- +goose StatementEnd
