create table tasks (
	id serial primary key,
	name text,
	description text,
	priority integer,
	due_date date,
	status integer,
	username text foreign key references users(username)
);
