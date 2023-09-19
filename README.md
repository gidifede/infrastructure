
# Infrastructure

This project will contain the CDK code to deploy the infrastructure for the logistic backbone.
The documentation of this project will be focused around CDK details and choices


# Taskfile commands

To install Taskfile:
https://taskfile.dev/installation/

`task init`
Initialize python dev environment

`task sync-services`
Update all submodule repos with remote

`task services-switch-branch -- [master|develop|branchname]`
Switch to the specified branch all submodule repos

`task update-graph`
Update the architecture image based on current CDK script

`task run-unit-tests`
Update all services and run only unit tests

`task run-integration-tests` or `task run-integration-tests -- [aws-profile]`
Update all services and run only integration tests

# Welcome to your CDK Python project!

This is a blank project for CDK development with Python.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

This project is set up like a standard Python project.  The initialization
process also creates a virtualenv within this project, stored under the `.venv`
directory.  To create the virtualenv it assumes that there is a `python3`
(or `python` for Windows) executable in your path with access to the `venv`
package. If for any reason the automatic creation of the virtualenv fails,
you can create the virtualenv manually.

To manually create a virtualenv on MacOS and Linux:

```
$ python3 -m venv .venv
```

After the init process completes and the virtualenv is created, you can use the following
step to activate your virtualenv.

```
$ source .venv/bin/activate
```

If you are a Windows platform, you would activate the virtualenv like this:

```
% .venv\Scripts\activate.bat
```

Once the virtualenv is activated, you can install the required dependencies.

```
$ pip install -r requirements.txt
```

At this point you can now synthesize the CloudFormation template for this code.

```
$ cdk synth
```

To add additional dependencies, for example other CDK libraries, just add
them to your `setup.py` file and rerun the `pip install -r requirements.txt`
command.

## Constructs

![Constructs](docs/constructs.png)
[ArchitetturaCDK.drawio](docs/ArchitetturaCDK.drawio)

## Useful commands

 * `cdk ls`          list all stacks in the app
 * `cdk synth`       emits the synthesized CloudFormation template
 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk docs`        open CDK documentation

Enjoy!
