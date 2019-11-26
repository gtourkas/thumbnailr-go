# Thumbnailr

This is a sample AWS serverless api for creating thumbnails from photos. It's built with Golang and 
[AWS SAM](https://github.com/awslabs/serverless-application-model) and its intented usage is to serve as a sandbox or an
educational resource.

## Getting Started

You need to follow the 
[AWS instructions](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) 
to install the required software and setup your AWS account.  

You also need to install [mage](https://magefile.org/) for building the lambdas and 
[delve](https://github.com/go-delve/delve) for locally debugging them. 

Note that I've developed this on a Mac OS and I'm using bash for a few, very simple, convenient shell scripts. Some of 
these scripts have a dependency to [jq](https://stedolan.github.io/jq/) and some others to [httpie](https://httpie.org/). 

## Build

The following snippets assume that you've opened a terminal and cd'ed to the the `host_lambda` directory.

### Build all Lambdas

```
mage -d ./mage build
```

which builds up to `max_concur_builds` lambdas concurrently. 
The value of `max_concur_builds` can be set in the `./mage/conf.ini` file.  

### Build a Single Lambda

```
TN_HANDLER=<LambdaResourceName> mage -d ./mage build
```

where LambdaResourceName is the resource name in the `template.yaml` file (e.g. `CheckCreationFunction`)

## Run Locally

For the lambdas triggered by the API Gateway, you need to run the following:

```
sam local start-api
```

then get the access token:

```
http POST "http://127.0.0.1:3000/token?grant_type=password&username=test&password=test"
```

and then invoke the lambda (e.g.):

```
http POST "http://127.0.0.1:3000/request_creation?width=100&length=100&photoID=buddha.jpg&format=PNG" authorization:"Bearer <ACCESS_TOKEN>"
```

The thumbnail creation lambda is triggered by an SNS event. To run this lambda you need to: 

1. edit the create_message.json
2. run `gen_create_event.sh` that generates the create_message.json
3. run `sam local invoke` as shown below

```
sam local invoke CreateFunction -e create_event.json
```  

Make sure the photstore S3 bucket has some photos. You can find a few in the `photostore` bucket.

## Debug Locally

The following snippets assume that you've opened a terminal and cd'ed to the the `host_lambda` directory.

For step-through debugging you need to build delve for the lambda OS and Architecture. This is done with:  

```
./build_delve.sh
```

Also you need to setup you IDE to run remote go debugging on port 5986.

These are one-off steps which you don't need to repeat. 

Then you need to build with env var `TN_DEBUG=true`. For all lambdas:

```
TN_DEBUG=true mage -d ./mage build
```

For a single lambda:

```
TN_DEBUG=true TN_HANDLER=<LambdaResourceName> mage -d ./mage build
```

where LambdaResourceName is the resource name in the `template.yaml` file (e.g. `CheckCreationFunction`)

Finally you need to follow the steps in the `Run Locally` section. Each `sam` command needs the following additions: 

```
-d 5986 --debugger-path ./delve --debug-args "-delveAPI=2"
``` 

## Deploy

Make sure you've built all handlers with debugging off (default option). With a terminal cd'ed to the the `host_lambda` 
directory run the following:

```
./deploy.sh
```

This script assumes an S3 Bucket named `sam-deployment-<AWS_ACCOUNT_ID>` where the packaged lambdas are uploaded.

Note that:
1. API GW HTTPS Url is printed. You need that Url to perform calls at the production environment.
2. ALB is not created (SAM does not provide this option yet). One needs to create it manually and wire it with the 
lambdas ending in '-alb'.  

## Run

The steps are identical to those of `Run Locally` but the local HTTP url needs to be replaced with the output HTTPS URL 
of the deployment script. 