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