CREATE TABLE ops_eml_number_successes
(
    id            BIGINT AUTO_INCREMENT PRIMARY KEY,
    cid           BIGINT       NOT NULL,
    from_address  VARCHAR(255) NOT NULL COMMENT '发件人邮箱',
    domain        VARCHAR(255) NOT NULL COMMENT '收件域',
    start_time    TIMESTAMP    NOT NULL COMMENT '区间开始时间',
    end_time      TIMESTAMP    NOT NULL COMMENT '区间结束时间',
    success_cnt   INT          NOT NULL COMMENT '成功发送数量',
    total_cnt     INT          NOT NULL COMMENT '总发送数量',
    temporary_err INT          NOT NULL DEFAULT 0 COMMENT '临时错误（可重试）',
    permanent_err INT          NOT NULL DEFAULT 0 COMMENT '永久错误（不可重试）',
    create_time   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP NOT NULL,
    update_time   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

alter table eml_tasks
    change tag user_tag varchar(64) not null comment '用于筛选用户';

# access_key_id 和 access_key_secret 需要加密存储，使用 encrypt_test.go文件加密
INSERT INTO eml_sys_users (cid, type, host, email_server, port, username, default_speed, password, access_key_id,
                           access_key_secret, from_address)
VALUES (3, 'register', 'mule', 'unknown', 465, 'mule', '1,2', 'mule', '', '', 'mule');

alter table eml_sys_users
    add default_speed varchar(16) default '16,1' not null comment '默认发件速度' after username;

alter table eml_templates
    add priority int default 100 not null comment '优先级（越小越好）' after template_name;

alter table eml_templates
    add max_retry int default 0 not null comment '最大重试次数' after priority;

