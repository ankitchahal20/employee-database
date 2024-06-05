# Simple Employee Database

This repository contains the source code for a Simple Employee Database built using Golang. The system is responsible for creating, deleting, updating and fetching of an employee record

## Prerequisites

Before running the Simple Employee Database, make sure you have the following prerequisites installed on your system:

- Go programming language (go1.21.4)
- PostgreSQL(14.8)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ankitchahal20/employee-database.git
   ```

2. Navigate to the project directory:

   ```bash
   cd employee-database
   ```

3. Install the required dependencies:

   ```bash
   go mod tidy
   ```

4. DB setup
    ```
    Use the scripts inside sql-scripts directory to create the tables in your db.
    ```
5. Defaults.toml
Add the values to defaults.toml and execute `go run main.go` from the cmd directory.

## APIs
There are five API's which this repo currently has.

Create Employee Record
```
curl -i -k -X POST \
   http://localhost:8080/v1/employees \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
"name":"Ankit Chahal",
"position": "Sr. Software Developer",
"salary": 12345.00
}'
```


Updating Employee Record

```
curl -i -k -X PUT \
  http://localhost:8080/v1/employees \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
  "id":"2",
  "name":"hell1o",
  "position": "some position",
  "salary": 123409.00
}'
```

Get Employee Record By ID

```
curl -i -k -X GET \
  http://localhost:8080/v1/employees/:id \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json"
```

Delete Employee Record

```
curl -i -k -X DELETE \
  http://localhost:8080/v1/employees/:id \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" 
```

## Project Structure

The project follows a standard Go project structure:

- `config/`: Configuration file for the application.
- `internal/`: Contains the internal packages and modules of the application.
  - `config/`: Global configuration which can be used anywhere in the application.
  - `constants/`: Contains constant values used throughout the application.
  - `db/`: Contains the database package for interacting with PostgreSQL.
  - `middleware`: Contains the logic to validate the incoming request
  - `models/`: Contains the data models used in the application.
  - `employeeerror`: Defines the errors in the application
  - `service/`: Contains the business logic and services of the application.
  - `server/`: Contains the server logic of the application.
  - `utils/`: Contains utility functions and helpers.
- `main.go`: Main entry point of the application.
- `README.md`: README.md contains the description for the employee-database.

## Contributing

Contributions to the Simple Employee Database are welcome. If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

The Simple Employee Database is open-source and released under the [MIT License](LICENSE).

## Contact

For any inquiries or questions, please contact:

- Ankit Chahal
- ankitchahal20@gmail.com

Feel free to reach out with any feedback or concerns.
