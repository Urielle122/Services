CREATE Table IF NOT EXISTS services(
    id VARCHAR(255) PRIMARY KEY,
    content TEXT,
    action VARCHAR(50) default 'NEXT',
    previous_content VARCHAR(255) default NULL,
    created_date datetime,
    last_modified_date datetime,
    title VARCHAR(20),
    type_action VARCHAR(45),
    type_action_name VARCHAR(45)
);