# Gogen Domain

Calling `gogen domain yourdomainname` will

* Provide all those files 

  ```
  ├── Dockerfile
  ├── README.md
  ├── config.json
  ├── config.sample.json
  ├── docker-compose.yml
  ├── domain_yourdomainname
  │   └── README.md
  ├── go.mod
  └── gitignore
  ```

* Generate unique random secret key for your application in `config.json`
* Copy the basic template for controller, gateway and crud