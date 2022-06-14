--- book_dir table
create table if not exists `book_dir` (
    `id` integer primary key autoincrement,
    `path` text not null,
    `create_time` date default current_date not null,
    `update_time` date default current_date not null
);