-- +goose Up
-- 新建外键
ALTER TABLE assistant_tools ADD CONSTRAINT assistant_tools_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);
ALTER TABLE assistant_tools ADD CONSTRAINT assistant_tools_tool_id_foreign FOREIGN KEY (tool_id) REFERENCES tools (id);
ALTER TABLE assistant_shares ADD CONSTRAINT  assistant_shares_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);
-- +goose Down
-- 删除外键
ALTER TABLE assistant_tools DROP FOREIGN KEY assistant_tools_assistant_id_foreign;
ALTER TABLE assistant_tools DROP FOREIGN KEY assistant_tools_tool_id_foreign;