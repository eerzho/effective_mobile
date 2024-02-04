create index idx_users_name on users using btree(name);
create index idx_users_surname on users using btree(surname);
create index idx_users_patronymic on users using btree(patronymic);
create index idx_users_age on users using btree(age);
create index idx_users_gender on users using btree(gender);
create index idx_users_country_id on users using btree(country_id);