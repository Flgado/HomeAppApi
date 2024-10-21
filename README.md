# Sam
This is a sample template for creating an API with two endpoints and an Amazon Cognito user pool that will protect the API. The two endpoints include:
1.**GetUserImageProfile**: Retrieves a profile image from an S3 bucket.
2. **PublishMessage**: Publishes a message to AWS IoT Core.

# Api Projection
- The API is secured by Amazon Cognito.
- Only users in the AdminGroup can send requests to the PublishMessage endpoint (which publishes messages to AWS IoT Core).
- Any authenticated user in the Cognito user pool can call the GetUserImageProfile endpoint to retrieve their profile image.

For detailed instructions on setting this up, refer to this [blog post](https://flgado.github.io/blog/).

## Project Structure:
```bash
.
├── api 
│   ├── go.mod
│   ├── go.sum
│   ├── lambda
│   │   ├── getUserImageProfile
│   │   │   ├── main.go
│   │   │   └── main_test.go
│   │   └── mqttPublisher
│   │       ├── main.go
│   │       └── main_test.go
│   ├── Makefile
│   ├── samconfig.toml
│   └── template.yaml
├── cognito
│   ├── go.mod
│   ├── go.sum
│   ├── lambda
│   │   └── addGroupScopeToIdToken
│   │       └── main.go
│   ├── Makefile
│   ├── samconfig.toml
│   └── template.yaml
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target 

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

The `sam build` command is wrapped inside of the `Makefile`. To execute this simply run
 
```shell
make
```

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://localhost:3000/hello`

**SAM CLI** is used to emulate both Lambda and API Gateway locally and uses our `template.yaml` to understand how to bootstrap this environment (runtime, where the source code is, etc.) - The following excerpt is what the CLI will read in order to initialize an API and its routes:


To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modifies IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

You can find your API Gateway Endpoint URL in the output values displayed after deployment.

### Testing

We use `testing` package that is built-in in Golang and you can simply run the following command to run our tests:

```shell
go test -v .
```



