[![.github/workflows/cicd.yml](https://github.com/tcwitte/go-azure-function/actions/workflows/cicd.yml/badge.svg)](https://github.com/tcwitte/go-azure-function/actions/workflows/cicd.yml)
# go-azure-function
This is a demo project that shows how an Azure Function can be written in Go and deployed.

The Github Actions workflow contains this functionality:
* on each pull request update or push, the code is compiled, deployed to Azure and "tested" (by calling wget to connect to the function)
* pull request builds get their own function app that is removed when the pull request is closed
* pushes to the main branch deploy to a test function app and then (after a manual approval configured in the environment) to a production function app

Many things are probably needed before using this in production, such as:
* authentication (if needed)
* use of Application Gateway or Front Door instead of direct access to the function
* a plan for the storage account and function app that provides the desired availability and performance
* use of a deployment slot in the production job for blue-green deployment

Later I intend to add more functionality, such as integration with a storage account or database. The "integration test" currently done with wget could be easily replaced by running a Postman collection, for example.
