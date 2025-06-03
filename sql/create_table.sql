# 业务用户表增加一个 tags 字段
alter table usr_users
    add tags varchar(1024) default '' not null comment '用户标签，可以用多个，存储形式:["a","b","c"]';


# 创建邮件推送系统相关的表


# 存储兑换码
create table if not exists email_redemption_codes
(
    id                bigint auto_increment
        primary key,
    cid               bigint                                not null comment '商户id',
    task_sub_id       bigint      default 0                 not null comment '邮件任务 ID',
    code              varchar(64)                           not null comment '兑换码',
    type              varchar(32) default 'credit'          not null comment '兑换码类别（credit）',
    amount            int         default 0                 not null comment '兑换值(100)',
    end_time          timestamp                             not null comment '失效时间',
    start_time        timestamp                             not null comment '生效时间',
    amount_valid_time timestamp                             null comment '兑换后有效期',
    status            int         default 0                 not null comment '是否已经使用（0：未使用，1：已使用）',
    uid               varchar(64)                           null comment '被哪个用户使用',
    used_time         timestamp                             null comment '使用时间',
    create_time       timestamp   default CURRENT_TIMESTAMP not null comment '创建时间',
    update_time       timestamp   default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
)
    comment '存储一些兑换码';

# 系统级用户关于邮件的配置信息
create table if not exists email_sys_users
(
    id           bigint auto_increment
        primary key,
    cid          bigint                                 not null comment '商户ID，关联商户',
    host         varchar(64)                            not null comment 'smtp 服务器',
    port         int                                    not null comment '端口',
    username     varchar(128)                           not null comment 'smtp服务用户名',
    password     text                                   not null comment '对称加密后的密码（授权码）',
    from_address varchar(64)                            not null comment '发件人地址',
    reply_to     varchar(255) default ''                not null comment '回信地址',
    event_track  text                                   null comment '埋点事件变量，以 json 形式存储。存储内容为 kv 键值对数组',
    update_time  timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    create_time  timestamp    default CURRENT_TIMESTAMP not null comment '创建时间'
)
    comment '系统级用户关于邮件的配置信息';

# 邮件标签表
create table if not exists email_tags
(
    id          bigint auto_increment
        primary key,
    cid         bigint                                  not null comment '关联商户 ID',
    tag         varchar(64)                             not null comment '标签',
    description varchar(1024) default ''                not null comment '标签描述',
    update_time timestamp     default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    create_time timestamp     default CURRENT_TIMESTAMP not null comment '创建时间'
);

# 邮件任务
create table if not exists email_task_subs
(
    id                 bigint auto_increment
        primary key,
    cid                bigint                                 not null comment '关联商户 ID',
    task_id            bigint                                 not null,
    from_email         varchar(255)                           not null comment '发件人地址',
    to_user_id         varchar(128) default ''                not null comment '收件人用户 ID',
    template_id        bigint       default 0                 not null,
    status             varchar(32)  default 'unsent'          not null comment '任务状态(unsent、processing、send failure、sent、opened、clicked)',
    version            int          default 0                 not null comment '乐观锁',
    fetch_time         timestamp    default CURRENT_TIMESTAMP not null comment '上次获取事件',
    status_description varchar(128) default ''                not null comment '状态描述',
    concrete_state     varchar(32)  default ''                not null comment '邮件详细的状态',
    update_time        timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    create_time        timestamp    default CURRENT_TIMESTAMP not null comment '创建时间',
    constraint email_task_subs_task_id_to_user_id_uindex
        unique (task_id, to_user_id)
);

create index email_task_subs_status_fetch_time_index
    on email_task_subs (status, fetch_time);

# 发件任务
create table if not exists email_tasks
(
    id           bigint auto_increment
        primary key,
    cid          bigint                                null comment '商户ID，关联商户',
    from_email   varchar(255)                          not null comment '发信人地址',
    tag          varchar(64)                           not null comment '邮件标签: ["a","b","c"]',
    template_id  bigint                                not null comment '模板ID',
    status       varchar(32) default 'pending'         not null comment '发件任务状态（pending、processing、finished、invalid）',
    process_info text                                  null comment '发件任务处理信息',
    version      int         default 0                 not null comment '乐观锁',
    fetch_time   timestamp                             null comment '上次获取时间',
    total_num    int                                   not null comment '发件总数 等价于 收件人数量',
    success_num  int         default 0                 not null comment '成功发送的数量',
    fail_num     int         default 0                 null comment '邮件任务失败的数量',
    invalid_num  int         default 0                 not null comment '收件人地址有误总数',
    pending_num  int         default 0                 not null comment '等待被发送的邮件数量',
    open_num     int         default 0                 not null comment '收到邮件后，打开邮件的行为',
    open_num_de  int         default 0                 not null comment '对”打开”行为基于同一封邮件做去重',
    click_num    int         default 0                 not null comment '收到邮件后，打开邮件后，点击了邮件中的URL',
    click_num_de int         default 0                 not null comment '对”点击”行为基于同一封邮件做去重',
    update_time  timestamp   default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    create_time  timestamp   default CURRENT_TIMESTAMP not null comment '创建时间'
);

create index email_tasks_status_fetch_time_index
    on email_tasks (status, fetch_time);

# 邮件模板
create table if not exists email_templates
(
    id            bigint auto_increment
        primary key,
    cid           bigint                                 not null comment '关联商户 ID',
    category      varchar(32)  default 'common'          not null comment '模板类型，目前有：普通（common）、推荐奖励（redemption）',
    content_type  varchar(32)                            not null comment '模板类型(text/HTML 或 text/plain)',
    template_name varchar(64)                            not null comment '模板名称',
    subject_zh    varchar(255)                           not null comment '邮件中文标题',
    subject_en    varchar(255) default ''                not null comment '邮件英文标题',
    from_name     varchar(64)                            not null comment '发送人名称',
    content_zh    text                                   not null comment '中文模板',
    content_en    text                                   null comment '英文模板',
    description   varchar(512) default ''                not null comment '模板描述（非必需）',
    update_time   timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间',
    create_time   timestamp    default CURRENT_TIMESTAMP not null comment '创建时间'
);

# 兑换码和邮件任务关联表
create table if not exists email_redemption_task_sub
(
    id          bigint auto_increment
        primary key,
    task_sub_id bigint        not null comment '邮件任务 ID',
    type        varchar(32)   not null comment '目标兑换码类型',
    amount      int default 0 not null comment '目标兑换码的兑换值'
)
    comment '记录兑换码信息和邮件任务关系';



