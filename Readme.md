Test task
=========

This repository contains solution of the test task given to me.

Description of the task is stored in _Description.txt_ file

Requirements:
* To launch this code you should have installed Go language compiler.
* You should have configured AWS CLI on the host you will launch the code
* Your AWS CLI should be able to login to the AWS. Perhaps you have to set this env variables:
    - AWS_SDK_LOAD_CONFIG=true
    - AWS_REGION=<region-name>

Preliminary steps:
 - Cd to the project directory
 - Hit the command _go get_
Command to launch: _go run main.go_

Or you can download a precompiled binary and launch it

Improvement ways
----------------
 - Make server calls not blocking
 - Use channels for error handling
 - Remove Debug output
 - Organize code better
 - Add configuration by:
    - Environment variables
    - Command line arguments ("flag")
    - Some configuration file
 - Deploy by RPM package
 - Sign package by GPG



*Note 1:* It is better to launch this code on GNU/Linux

*Note 2:* If you use Mac OS, AWS EC2 instances will not be able to execute the deployed binary, because they are Linux not Mac OS