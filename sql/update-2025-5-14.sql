alter table eml_sys_users
    add access_key_id varchar(255) default '' not null after password;

alter table eml_sys_users
    add access_key_secret varchar(255) default '' not null after access_key_id;

alter table eml_sys_users
    add smtp_server varchar(64) default 'ali' not null after host;

