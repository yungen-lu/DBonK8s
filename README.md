# DBonK8s

>  A all-in-one tool for deploying different database with ease

## Motivation

The goal of this project is to provide a easy way for user to spin up a database for development or testing(mostly testing), it also provides  simple user management for admins.

## Overview

DBonK8s is a linebot running on kubernetes, user can interact with the bot to access and manage it's own databases(or others for admin).

## Features

* Spin up different types of database
* Config database default username and password
* Manage databse and view its info

## Available Commands

* `config [-upt]` 
* `list [-an]`
* `info [-dn]`
* `stop [-dn]`
* `create [-dtn]`
* `back`

## Deploy

### Prerequisite

* gcloud account
* terrafrom

### Enviroment setup

* setup a project in gcloud
* modify `project_id` and `region` in `terraform/terraform.tfvars`

### Deploy with terraform

```shell
terraform apply
```



## Finite State Machine

<img src="./assets/fsm.svg" alt="graphviz" style="zoom:50%;" />