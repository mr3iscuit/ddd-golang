# Create domain subfolders
mkdir -p domain/model domain/repository domain/service domain/event

# Move domain files into appropriate folders
mv domain/todo.go domain/model/todo.go
mv domain/repository.go domain/repository/todo_repository.go

# Create an empty domain service and event files
touch domain/service/todo_domain_service.go
touch domain/event/todo_completed.go

# Create application subfolders
mkdir -p application/command application/query application/service application/dto

# Move existing application files to correct folders and rename as requested
mv application/command.go application/command/add_todo_command.go
mv application/port.go application/query/list_todos_query.go
mv application/service.go application/service/todo_service.go
mv application/view.go application/dto/todo_view.go

# Create infrastructure subfolders
mkdir -p infrastructure/repository

# Move infrastructure file and rename
mv infrastructure/memmory_todo_repo.go infrastructure/repository/todo_repository_memory.go

# Create interface subfolders
mkdir -p interface/cli interface/http

# Move interface CLI file
mv interface/cli.go interface/cli/cli.go

# Create placeholder for HTTP handler
touch interface/http/todo_handler.go

# Create shared folder and placeholder
mkdir -p shared
touch shared/errors.go

# Move main.go to root folder (if main directory exists, merge its content)
mv main/main.go ./

# Remove now empty main folder
rmdir main
