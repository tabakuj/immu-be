# Notes
> This application is build using Golang,Gin and ImmuDb vault. I have structured this the same as I would structure i real application.
> <br/> So it might be a bit over-complicated for this task but pls keep in mind the reason.
> <br/> In addition I didn't add the pagination as it would be a little complicated on the front-end part
> ( if you guys want to see that on this please let me know).
> <br/>
> I have taken configuration and secrets from the same place for simplicity.
---

### 1. Prerequisites, Run test and Build image
I assume that Golang and Docker is installed on the machine you are trying to run this from.

* To run test please use the follow:
    ```shell
        #cd path to  solution
        go test -v ./...
    ```
* We might need to update the configuration setting for this <br>
    `config/default.yaml` is the configuration file.
    ```
    ImmuDbUrl: ""
    ImmuDbApiKey: ""
    ImmudbSearchUrl: ""
  ```
* Build the image<br>
    We will carry the config file we have on this solution for simplicity and use that.<br>
    In this step we will create an image with name immudb-docker-img (pls do not change it since it's used on the docker file as well)
    ```shell
        #cd path to  solution
        docker build  -t immudb-docker-img --build-arg FILEDIR="./config" --build-arg FILENAME="default.yaml" --build-arg TARGETFILENAME="default.yaml" .
    ```

### 2. Run Docker
We have to ways to run this either standalone or by using the docker compose.<br>
In this step, we will run the image built in the previous step.<br>
The application will run on port 8080, so we need to expose this port. Here, I have used port 8080.<br>


* Standalone 
<br>If you run this command, a container named "immuapi" will be created and running.
    ```shell
    # to remove the container 
    #docker rm -f immuapi
    docker run --name immuapi -p 8080:8080 -it immudb-docker-img        
    ```
  * with doccker compose <br>
    here the fe docker and be docker will run on 8081 and 8080 ports.
   ```shell
    docker compose up -d    
    ```
        
    >   **Note** <br>
    The fe image need to be created for the docker compose to work.


### 3. Endpoints 
There are 3 endpoints. all this information can be checked in [swagger](http://localhost:8080/swagger/index.html) as well but i wanted to include here also. 
* Create <br>
    here you can create a new account-info by providing a sample like the following.
    It will either return 201,400, or 500 if something unexpected happened.<br>
    sample call: 
    ```
  curl --location 'http://localhost:8080/v1/api/account-info' \
    --header 'Content-Type: application/json' \
    --data '{
    "account_name": "John Doe 22",
    "iban": "GB82WEST12345698765432",
    "address": "1234 Elm Street, Springfield, USA",
    "amount": 1500.75,
    "type": 1
    }'
  ```
* GetAll <br>
  this endpoint take 2 optional parameter for page and pageSize.<br>
    sample call:
    ```
  curl --location --request GET 'http://localhost:8080/v1/api/account-info?page=1&pageSize=2'
  ```
* GetByID <br>
 this endpoint search on the volt for the id requested.<br>
  sample call:
    ```
  curl --location 'http://localhost:8080/v1/api/account-info/3'
  ```

### 4. Code structure
The structure of the code is as follows: <br>
``` 
     Main Folder
        cmd --> main reside here, ususally this might be a command for example
        config --> the config for the application resides here              
        internal --> all the code resides on their respective packages
```


> **Note** <br>
> I have added some test on this application, but I don't consider this app covered in tests(for example there are not test for the handlers).

> **Note** <br>
> On this app i have logged the error directly in some places, usually I prefer that only the caller to log the error.
