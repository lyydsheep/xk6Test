alter table eml_task_subs
    add retry int default 0 not null comment '重试次数' after status;

alter table eml_task_subs
    add type varchar(36) default 'market' not null comment '邮件任务类型' after status;

alter table eml_tasks
    add type varchar(36) default 'market' not null comment '发件任务类型' after status;

create table eml_task_config
(
    id        bigint auto_increment
        primary key,
    cid       bigint                       not null,
    category  varchar(36) default 'email'  not null comment '任务类型(大分类）',
    type      varchar(36) default 'market' not null comment '任务类型（小分类）',
    priority  int         default 100      not null comment '值越小，优先级越高',
    max_retry int         default 3        not null comment '最大重试次数'
)
    comment '任务配置表';

alter table eml_sys_users
    drop column smtp_server;

INSERT INTO eml_task_config (id, cid, category, type, priority, max_retry) VALUES (1, 2, 'email', 'market', 100, 3);
INSERT INTO eml_task_config (id, cid, category, type, priority, max_retry) VALUES (2, 2, 'email', 'redemption', 10, 1000);

alter table eml_task_subs
    add priority int default 100 not null comment '数值越小，优先级越高' after type;


create table ops_eml_result
(
    id                   bigint auto_increment
        primary key,
    account_name         varchar(64)  default ''                not null comment '发件账号',
    error_classification varchar(64)  default ''                not null comment '错误类型',
    sent_time            timestamp                              not null comment '发送时间',
    message              varchar(255) default ''                not null comment '详情信息',
    status               int          default 0                 not null comment '投递结果',
    subject              varchar(255) default ''                not null comment '主题',
    to_address           varchar(255) default ''                not null comment '收件地址',
    create_time          timestamp    default CURRENT_TIMESTAMP not null,
    update_time          timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    constraint ops_eml_result_sent_time_to_address_account_name_uindex
        unique (sent_time, to_address, account_name)
);

alter table eml_sys_users
    add email_server varchar(32) default 'ali' not null comment '服务商' after host;

UPDATE eml_sys_users SET email_server = '' WHERE host = 'smtp.qq.com';

alter table ops_eml_result
    modify message varchar(2048) default '' not null comment '详情信息';

alter table eml_task_config
    add create_time timestamp default CURRENT_TIMESTAMP not null;

alter table eml_task_config
    add update_time timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP;

alter table eml_sys_users
    add start_time timestamp default '2025-05-09 00:00:00' not null comment '爬取任务起始时间' after password;