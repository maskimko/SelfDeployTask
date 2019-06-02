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

What should be done else:
-------------------------
 - Add InstanceProfile with attached policies and roles
   to the launched instances to allow deployed binary to act
 - Fix region guessing inside the EC2 instance (Use get http call to the magic endpoint)

Additional information:
-----------------------
 - Application listen to the TCP socket on port 1989
 - Stop command can be sent like this: _echo   "stop" | nc -4 localhost 1989_
 - Move command can be sent like this: _echo  -n "moveto 'eu-west-1'" | nc -4 localhost 1989_

Improvement ways
----------------
 - Make server calls not blocking
 - Use channels for error handling
 - Remove Debug output
 - Organize code better
 - Do not hardcode listen port and some AWS tag names
 - Add configuration by:
    - Environment variables
    - Command line arguments ("flag")
    - Some configuration file
 - Deploy by RPM package
 - Sign package by GPG


 * Note 0: Timings:
    - Spinning up an VPC takes about 5 minutes (Mainly because of waiting for EC2 instances to bootstrap)
    - Shut down the region takes about 2 minutes
 * Note 1: It is better to launch this code on GNU/Linux
 * Note 2: If you use Mac OS, AWS EC2 instances will not be able to execute the deployed binary, because they are Linux not Mac OS
 * Note 3: I temporary disable start of deployed binary
