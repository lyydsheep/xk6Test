# 2025 年 5 月 13 日

alter table eml_task_subs
    add code_smtp varchar(8) default '' not null comment 'smtp 服务器返回码' after fetch_time;

alter table eml_task_subs
    drop column status_description;

alter table eml_task_subs
    add code_description varchar(2048) default '' not null comment '返回码描述' after code_smtp;

alter table eml_task_subs
    drop column concrete_state;

alter table eml_task_subs
    add sent_time timestamp null comment '发件时间' after code_description;

alter table eml_task_subs
    add email_domain varchar(64) default '' not null comment '收件域' after sent_time;


create table eml_domain_credibility
(
    id                  bigint auto_increment primary key,
    cid                 bigint                                 not null,
    from_address        varchar(255)                           not null comment '发件地址',
    domain              varchar(255) default ''                not null comment '收件域名，@后面的东西，不包括@',
    success_rate_day    DECIMAL(5, 2)                          null comment '每日成功率',
    success_rate_hour   DECIMAL(5, 2)                          null comment '每小时成功率',
    success_rate_second DECIMAL(5, 2)                          null comment '每分钟成功率',
    speed               varchar(8)   default '16,1'            not null comment '发件速度 xx秒 xx 封 ---> 8秒 1 封：8,1    1秒 3 封：1,3',
    create_time         timestamp    default CURRENT_TIMESTAMP not null,
    update_time         timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
);

alter table eml_domain_credibility
    add last_time timestamp default '1970-01-01 00:00:01' not null comment '上一次向收件域发送邮件的时间' after domain;