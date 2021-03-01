# Webapp
Designed for c++ engineer developping web applications rapidly

![SystemDesign](https://github.com/zanlichard/webapp/raw/master/docs/SystemDesign.png)


# Index-Page
- [Webapp](#webapp)
- [Index-Page](#index-page)
- [Main Directory Guide](#main-directory-guide)
  - [Configure](#configure)
  - [Framework Entry](#framework-entry)
    - [Appengine](#Appengine)
    - [Appframework](#Appframework)
    - [Appframeworkboot](#Appframeworkboot)
    - [Appinterface](#Appinterface)
  - [Application Develop](#application-develop)
    - [Router](#router)
    - [Service](#service)
    - [Errors](#errors)
    - [Dao](#dao)
    - [Model](#model)
  - [Monitoring](#monitoring)
  - [Storage Middleware](#storage-middleware)
  - [Message Queue Middleware](#message-queue-middleware)
  - [Other](#other)
    - [Middleware](#middleware)
    - [Toolkit](#toolkit)
    - [Internal](#internal)
  - [Use Restrictions](#use-restrictions)
    - [Logging](#logging)
    - [ErrorCode](#errorcode)
    - [Monitoring Add](#monitoring-add)
  - [DeployMent](#deployment)
    - [Recommand Directory](#recommand-directory)
    - [Extend Directory](#extend-directory)

# Main Directory Guide

## Configure
Application main config is toml format

## Framework Entry
### Appengine
Designed for setup framework instance

### Appframework
Define framework level errcode 401,400 for example<br>
Define service local config<br>
Define service dependent config<br>
Define framework global variables<br>

### Appframeworkboot
Setup associate config such as crontab task for delete expire log<br>
Init framework global varibles<br>

### Appinterface
Define application api interface,such as json supported protocol object

## Application Develop

### Router
Define application interface url

### Service
Define business logic

### Errors
Define business error code and error message<br>
Support local call statck<br>

### Dao
Define data access according to mysql redis mongo driver

### Model
Define data object,such as mysql table object 


## Monitoring
When app runs,it will push api call stat,errcode stat,delay stat,in and out flow stat info into stat log file at regular intervals

## Storage Middleware
Support mysql based on gorm<br>
Support redis based on redisgoe<br>
Support mongo<br>

## Message Queue Middleware
Just support rabbitmq


## Other
### Middleware
Support token based on jwt<br>
Support interface sign<br>

### Toolkit
Define custom function 

### Internal
Support promethus metrics scan<br>
Support pprof performance data scan,such as goroutines,threads,heap info<br>

## Use Restrictions

### Logging
Dao layer suggest printing no log<br>
Api layer print detail log about http<br>
Service layer print business log<br>

### ErrorCode
Dao layer return base errcode<br>
Service layer return business define errcode<br>
Api layer may return standard error base on http and business errcode<br>

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
├── bin<br>
│   └── webapp<br>
├── etc<br>
│   └── config.toml<br>
├── frameworklog<br>
│   ├── access<br>
│   │   └── webapp<br>
│   │       ├── log.2021-02-22<br>
│   │       └── log.2021-02-23<br>
│   ├── business<br>
│   │   └── webapp<br>
│   │       ├── log.2021-02-22<br>
│   │       └── log.2021-02-23<br>
│   └── err<br>
│       └── webapp<br>
│           ├── log.2021-02-22<br>
│           └── log.2021-02-23<br>
├── log<br>
│   ├── app.2021-02-23.001.log<br>
│   └── app.log<br>
├── logs<br>
│   └── nohup.err<br>
├── stat<br>
│   └── webapp_stat.log<br>
└── tools<br>
│   └── op<br>
│   │     ├── p.sh<br>
│   │     ├── start.sh<br>
│   │     └── stop.sh<br>

