UPDATE todo_models
SET title = 'updated curl!!', completed = 1
WHERE id = 2
RETURNING id;
