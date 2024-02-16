# BANK STATEMENT PARSER

## Local environment set up:
1. Create/Update .env file: 
   - The service expects a *.env* file to be present in the root level directory of the project. 
   - Following are the environment variables that help to configure the behaviour of the service.

    | Variable Name           | Usage                                              | Mode     | Default Value          | Sample Value           |
    | ----------------------- | -------------------------------------------------- | -------- | ---------------------- | ---------------------- |
    | APP_MODE                | service can be set up in two mode: CMD or HTTP     |          |                        | CMD                    |
    | FILE_COLUMN_SEPARATOR   | character used to separate columns in file         | CMD/HTTP | ,                      | ,                      |
    | FILE_HAS_HEADER         | flag to specify if file has a header row           | CMD/HTTP | false                  | true                   |
    | PAYMENT_REFERENCE_REGEX | regular expression to check the payment reference  | CMD/HTTP | PAY[0-9]{6}[a-zA-Z]{2} | PAY[0-9]{6}[a-zA-Z]{2} |
    | DECIMAL_PRECISION       | decimal precision for credit/debit amounts in file | CMD/HTTP | 2                      | 3                      |
    | SERVER_ADDR             | http server address {HOST:PORT}                    | HTTP     | localhost:8080         | :9001                  |
    | ENABLE_TLS              | flag to enable HTTPS                               | HTTP     | false                  | true                   |
    | SSL_CRT_PATH            | path to SSL certificate file                       | HTTP     |                        | cert.pem               |
    | SSL_KEY_PATH            | path to SSL key file                               | HTTP     |                        | key.pem                |
 
2. Start the service: 
   - Run the below command to start the service. 
        ```
        go run main.go
        ```
   - If **APP_MODE** configured as *CMD*, you should see command line instructions to use the service.
   - If **APP_MODE** configured as *HTTP*, http server will be started on host and port specified by *SERVER_ADDR* variable. You can use any http client to use the API. The API definition is available on swagger endpoint **http://{SERVER_ADDR}/docs/index.html**


## Generate/View Swagger for HTTP mode:

   - Install the swag binary using the below command: **Required for generating Swagger files**
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
