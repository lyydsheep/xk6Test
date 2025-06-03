alter table ops_eml_result
    add domain varchar(64) default '' not null comment '收件域' after account_name;