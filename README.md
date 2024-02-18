# BANK STATEMENT PARSER
This service parses the bank statement file to deduce the total value of payments made on a given date for each currency. All lines of the bank statement that represent a payment should contain a payment reference in one of it's narratives (refer to sample file *data/statement.csv*). Refer below section for usage and other details.

## Usage

1. Create/Update .env file: 
   - The service expects a *.env* file to be present in the root level directory of the project if required configuration is not set as environment variables. 
   - Following are the environment variables that help to configure the behaviour of the service.

    | Variable Name           | Usage                                               | Mode | Default Value          | Sample Value           |
    | ----------------------- | --------------------------------------------------- | ---- | ---------------------- | ---------------------- |
    | APP_MODE                | service can be set up in two mode: CMD or HTTP      | BOTH | CMD                    | CMD                    |
    | FILE_COLUMN_SEPARATOR   | character used to separate columns in file          | BOTH | ,                      | ,                      |
    | FILE_DATE_FORMAT        | format (golang date format) for date column in file | BOTH | 02/01/2006             | 2006-01-02             |
    | PAYMENT_REFERENCE_REGEX | regular expression to check the payment reference   | BOTH | PAY[0-9]{6}[a-zA-Z]{2} | PAY[0-9]{6}[a-zA-Z]{2} |
    | DECIMAL_PRECISION       | decimal precision for credit/debit amounts in file  | BOTH | 2                      | 3                      |
    | SERVER_ADDR             | http server address {HOST:PORT}                     | HTTP | localhost:8080         | :9001                  |
 
2. Start the service: 
   - Build the executable using below command. 
        ```
        go env -w GOOS=linux
        go build -o ./bin/main main.go
        ```
     You should get main or main.exe file based on value of GOOS go env variable

     Already built linux and windows executables can be found in **bin** folder
   - Run the executable using below command.
     - If **APP_MODE** configured as *CMD*, following command line flags are available to use the service.  
          | Flag Name     | Description                                                       |
          | ------------- | ----------------------------------------------------------------- |
          | date          | filter date format: DD/MM/YYYY  (if not passed all rows are used) |
          | file_path     | path to input csv file                                            |
          | out_file_path | if provided the result will be stored in file (optional)          |

     - Use the below command to start the service. 
        ```
        go run main.go --date <Date e.g. 06/03/2011> --file_path <path to input csv>
        ```
        OR
        ```
        ./main --date <Date e.g. 06/03/2011> --file_path <path to input csv>
        ```
  
   - If **APP_MODE** configured as *HTTP*, http server will be started on host and port specified by *SERVER_ADDR* variable. 
     - Use the below command to start the service. 
        ```
        go run main.go
        ```
        OR
        ```
        ./main
        ```
     - You can use any http client to use the API. 
     - The API definition is available on swagger endpoint **http://{SERVER_ADDR}/docs/index.html**

## Containerization
   - The dockerfile present in the root directory can be used to create docker image of the service.
   - Use the below command to create docker image
     ```
     docker build -t statement_processor .
     ```
     You can find the docker image with latest code at [vinitondocker/statement_processor](https://hub.docker.com/repository/docker/vinitondocker/statement_processor/general)
   - The usage remains same as mentioned above. Below mentioned are some sample commands
     
     **APP_MODE** = *HTTP*
     ```
     docker run -it --env-file .env --name app -p 9001:9001 statement_processor
     ```
          
     **APP_MODE** = *CMD*
     ```
     docker run -it --env-file .env --v ${pwd}/data:/app/data --name app statement_processor --date 06/03/2011 --file_path data/statement.csv
     ```


## Generate/View Swagger for HTTP mode:

   - Install the swag binary using the below command: **required for generating Swagger files**
        ```
        go install github.com/swaggo/swag/cmd/swag@latest
        ``` 
    
   - To generate run the swagger generation command from root folder of the project 
        ```
        swag init
        ```
   - The above command should generate a docs folder if it doesnt already exists.
   - Three files, namely *docs.go, swagger.json* and *swagger.yaml* will be created in the *docs* folder.
   - After every update to swagger related comments in code it is required to run *swag init* command to update the documentation.
   - Swagger document will be available on **http://{SERVER_ADDR}/docs/index.html**
