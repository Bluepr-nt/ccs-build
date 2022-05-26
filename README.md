# build-action

Backbone implementation of the build task as per the [CI-CD-Standard](https://github.com/13013SwagR/CI-CD-Standard)

## Features
```
Feature: Execute build command in environment

  Background:
    Given a build environment 
    Given a build command

  Scenario Outline: various build command & environment
    Given a <environment_type> build environment 
    Given a <command_type> build command
    When the build is triggered
    Then the build command is executed in the environment
    And outputs
    * the result code (integer)
    * the execution log artifact(url)
    * the pipeline/build artifact(url)
    
    Examples:
      | environment_type | command_type |
      | Docker-in-Docker |       Docker |
      |     Docker image |          tar |
      |     Docker image |           go |

```