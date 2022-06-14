---book_list
create table if not exists `book_list` (
    `id` integer primary key autoincrement,
    `name` text not null unique,
    `md5` text not null unique,
    `create_time` date default current_date not null,
    `update_time` date default current_date not null
);