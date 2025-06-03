update eml_task_subs set `type` = 'redemption' where template_id = 2;
update eml_task_subs set `type` = 'market' where template_id = 1;

update eml_task_subs set priority = 10 where `type` = 'redemption';
update eml_task_subs set priority = 100 where `type` = 'market';