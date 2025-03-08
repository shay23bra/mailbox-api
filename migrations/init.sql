-- Create departments table
CREATE TABLE IF NOT EXISTS departments (
    department_id INT PRIMARY KEY,
    department_name VARCHAR(100) NOT NULL
);

-- Create mailboxes table
CREATE TABLE IF NOT EXISTS mailboxes (
    mailbox_identifier VARCHAR(100) PRIMARY KEY,
    user_full_name VARCHAR(100) NOT NULL,
    job_title VARCHAR(100) NOT NULL,
    department_id INT NOT NULL,
    manager_mailbox_identifier VARCHAR(100),
    org_depth INT NOT NULL DEFAULT 0,
    sub_org_size INT NOT NULL DEFAULT 0,
    FOREIGN KEY (department_id) REFERENCES departments(department_id),
    FOREIGN KEY (manager_mailbox_identifier) REFERENCES mailboxes(mailbox_identifier)
);

-- Create indexes
CREATE INDEX idx_mailboxes_department_id ON mailboxes(department_id);
CREATE INDEX idx_mailboxes_manager_id ON mailboxes(manager_mailbox_identifier);
CREATE INDEX idx_mailboxes_org_depth ON mailboxes(org_depth);
CREATE INDEX idx_mailboxes_sub_org_size ON mailboxes(sub_org_size);