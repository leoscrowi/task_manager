# Task Manager API
_______________

## Как запустить task-service:
- Есть поддержка multi-stage builds
```bash
docker compose -f 'docker-compose.yml' up -d --build 
```
________________

## Base URL
```
http://localhost:8080/task
```

## API Endpoints

### 1. Create Task
- **Method**: `POST`
- **Endpoint**: `/task`
- **Description**: Create and save a new task.

#### Request Body:
```json
{
    "user_id": "b063de04-6fd7-41cd-8f4c-8d113e786be8",
    "title": "Sample Task",
    "description": "This is a sample task description.",
    "repeat_task": "DAILY",
    "parent_task_id": "b063de04-6fd7-41cd-8f4c-8d113e786be8"
}
```
#### Responses:
- 201 Created: Task created successfully.
- 400 Bad Request: Invalid request parameters.
- 500 Internal Server Error: Server error during task creation.


### 2. Update Task
- **Method**: `PATCH`
- **Endpoint**: `/task`
- **Description**: Update details of an existing task by UUID.

#### Request Body:
```json
{
    "title": "Updated Task Title",
    "description": "Updated task description.",
    "repeat_task": "WEEKLY",
    "task_status": "IN_PROGRESS"
}
```
#### Responses:
- **200 OK**: Task updated successfully.
- **400 Bad Request**: Invalid request parameters.
- **500 Internal Server Error**: Server error during task update.

### 3. Delete Task
- **Method**: `DELETE`
- **Endpoint**: `/task`
- **Description**: Delete task by UUID.

#### Request Body:
```json
{
  "title": "Sample Task",
  "description": "This is a sample task description.",
  "repeat_task": "DAILY" // Options: DAILY, WEEKLY, MONTHLY, YEARLY, NEVER
}

```
#### Responses:
- **200 OK**: Task deleted successfully.
- **400 Bad Request**: Invalid request parameters.
- **500 Internal Server Error**: Server error during task deleting.


### 4. Create Task
- **Method**: `POST`
- **Endpoint**: `/task`
- **Description**: Creates and saves a new task.

#### Request Body:
```json
{
    "id": "b063de04-6fd7-41cd-8f4c-8d113e786be8",
}
```
#### Responses:
- **200 OK**: Task deleted successfully.
- **400 Bad Request**: Invalid request parameters.
- **500 Internal Server Error**: Server error during task deleting.
