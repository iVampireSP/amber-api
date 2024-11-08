-- +goose Up

CREATE TABLE scene_prompts (
  id serial primary key,
  label varchar(255) not null,
  prompt text not null,
  assistant_id bigint unsigned not null,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  foreign key (assistant_id) references assistants(id) on delete cascade,
  index scene_prompts_assistant_id_idx (assistant_id),
  index scene_prompts_label_idx (label)
);

-- +goose Down
drop table scene_prompts;