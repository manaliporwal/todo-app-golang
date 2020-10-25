This repository is an implementation of simple ToDo list application in Go.
User accounts and ToDo tasks are stored in a Postgres DB instance.

To run the application:

1) Create a postgres DB, users table using the users.sql and tasks table using tasks.sql
2) Start the application with the commands: go build and ./todo-app-golang
3) Test the application with the requests in the test.http file

PostgreSQL Tables:
------------------
create table users (
          username text primary key,
          password text
);

create table tasks (
	id serial primary key, 
	name text, 
	description text, 
	priority integer,
	due_date date, 
	status integer,
	username text foreign key references users(username)
);


API Enpoints:
-------------
POST '/signup'
This API is used for create a new account for the user. 
This is a POST request, the JSON fields for username and password must be provided in the input request body.



POST '/signin'
This API is used for user login.
This is a POST request, the JSON fields for username and password must be provided in the input request body.
Once the login is success, there is a session cookie generated which is stored in the response header.
User has to present this cookie for accessing the future ToDo list requests.

 
POST '/addTask'
This API adds a new ToDo task for the user.
This is a POST request, the JSON fields for task name, description, priority, due_date, status must be provided in the input request body.
For authentication of the user, the session cookie has to be presented in the request header.


GET '/getTasks'
This API returns a list of all the ToDo tasks created by the user.
For authentication of the user, the session cookie has to be presented in the request header.


PUT '/updateTask/{id}'
This API updates a ToDo task for the user.
The {id} field for a ToDo task is returned in the '/get/Tasks' API.
This is a PUT request, the JSON fields for task name, description, priority, due_date, status must be provided in the input request body.
For authentication of the user, the session cookie has to be presented in the request header.



PUT '/completeTask/{id}'
This API marks a ToDo task as Completed for the user.
The {id} field for a ToDo task is returned in the '/get/Tasks' API.



DELETE '/deleteTask/{id}'
This API deletes a ToDo task for the user.
The {id} field for a ToDo task is returned in the '/get/Tasks' API.

