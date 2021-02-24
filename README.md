# Webapp
Designed for c++ engineer developping web applications rapidly


# Index-Page
[toc]

# Main Directory Guide

## Configure
Application main config is toml format

## Framework Entry
### appengine
Designed for setup framework instance

### appframework
Define framework level errcode 401,400 for example
Define service local config
Define service dependent config
Define framework global variables

### appframeworkboot
Setup associate config such as crontab task for delete expire log
Init framework global varibles

### appinterface
Define application api interface,such as json supported protocol object

## Application Develop

### Router
Define application interface url

### Service
Define business logic

### Errors
Define business error code and error message
Support local call statck

### Dao
Define data access according to mysql redis mongo driver

### Model
Define data object,such as mysql table object 


## Monitoring
When app runs,it will push api call stat,errcode stat,delay stat,in and out flow stat info into stat log file at regular intervals

## Storage Middleware
Support mysql based on gorm
Support redis based on redisgoe 
Support mongo

## Message Queue Middleware
Just support rabbitmq


## Other
### Middleware
Support token based on jwt
Support interface sign

### Toolkit
Define custom function 

### Internal
Support promethus metrics scan
Support pprof performance data scan,such goroutines,threads,heap info

## Use Restrictions

### Logging
Dao layer suggest printing no log
Api layer print detail log about http
Service layer print business log

### ErrorCode
Dao layer return base errcode
Service layer return business define errcode
Api layer may return standard error base on http and business errcode

### Monitoring Add
Monitor keyword define in api layer,init at router layer


## DeployMent
### Recommand Directory
├── bin<br>
├── etc<br>
├── frameworklog<br>
├── log<br>
├── logs<br>
├── stat<br>
└── tools<br>

### Extend Directory
├── bin
│   └── webapp
├── etc
│   └── config.toml
├── frameworklog
│   ├── access
│   │   └── webapp
│   │       ├── log.2021-02-22
│   │       └── log.2021-02-23
│   ├── business
│   │   └── webapp
│   │       ├── log.2021-02-22
│   │       └── log.2021-02-23
│   └── err
│       └── webapp
│           ├── log.2021-02-22
│           └── log.2021-02-23
├── log
│   ├── app.2021-02-23.001.log
│   └── app.log
├── logs
│   └── nohup.err
├── stat
│   └── webapp_stat.log
└── tools
    └── op
        ├── p.sh
        ├── start.sh
        └── stop.sh

