create index eml_task_subs_status_retry_priority_index
    on eml_task_subs (status, retry, priority);

alter table eml_domain_credibility
    change last_time last_sent_time timestamp default '1970-01-01 00:00:01' not null comment '上一次向收件域发送邮件的时间';

alter table eml_sys_users
    change start_time fetch_start_time timestamp default '2025-05-09 00:00:00' not null comment '爬取任务起始时间';

