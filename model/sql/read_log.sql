---read_log
create table if not exists `read_log` (
    `id` integer primary key autoincrement,
    `book_id` integer unique not null,
    `line_count` integer not null,
    `create_time` date default current_date not null,
    `update_time` date default current_date not null
);