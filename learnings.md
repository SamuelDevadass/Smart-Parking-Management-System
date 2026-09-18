my-go-project/
├── main.go            # Entry point (initializes DB, router, starts server)
├── handlers/          # Equivalent to FastAPI routers (HTTP layer)
├── services/          # Equivalent to FastAPI services (Business logic)
├── repositories/      # Equivalent to database/backend calls (Data access)
└── models/            # Equivalent to Pydantic/SQLAlchemy models (Structs)
![architecture](image-1.png)
![flow path](image.png)  

go mod init is like setting up your project environment or running pip init / initializing poetry.

go.mod plays the exact same role as requirements.txt or a lock file, listing all the external packages your project depends on and their versions.

go get is essentially pip install.

go slice is like python list but fixed element types unlike python list

%v stands for "value" and is used with functions like fmt.Printf or fmt.Sprintf to inject variables into strings.

(Bonus tip: Go also has %+v to print struct fields with their names, and %T to print the data type!)

map is like a dictionary but unordered so looping on a map is not possible

In Go, %q stands for quoted string.

When you use %q in formatting functions like fmt.Sprintf or testing functions like t.Errorf, it takes whatever string or character you pass it and automatically wraps it in double quotes, safely escaping any special characters inside it.
For example:

If msg is Hello, Gladys, printing it with %v gives: Hello, Gladys

Printing it with %q gives: "Hello, Gladys"
Also, test functions take a pointer to the testing package's testing.T type as a parameter. You use this parameter's methods for reporting and logging from your test.

15-09-2026

While creating a go file, if its a standalone file that must run by itself or is the main entry point it must have the line package main

the main func is the main entry exit point of that file

use the joho/godotenv to to use .env files or os.Getenv if loading directly from OS

run go mod init to create the go.mod file 
after writing, run go mod tidy to update the mod file
to compile :
1) run go build -o executable_name code_file.go
run ./executable_name.exe or just ./executable_name
2) run go build
this creates an executable of name same as the direectory 
run ./executable_name.exe or just executable_name
the go env looks for 1st @ so if pw contains @ write as %40 alwys use the codes for the special characters in the passwords


In Go, maps are reference types. When you declare single_row := make(map[string]string) outside of a loop and keep appending it, you aren't adding new items—you are appending references to the exact same map. By the time the loop finishes, every single item in your list will look like the very last row fetched.