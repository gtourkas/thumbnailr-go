# Thumbnailr

## About




## Getting Started

You need to follow the [AWS instructions](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) to install the required software and setup your AWS account.  

You also need to install [mage](https://magefile.org/) for building the lambdas and [delve](https://github.com/go-delve/delve) for locally debugging them. 

Note that I've developed this in `Mac OS` and I'm using `bash` for a few, very simple, convenient shell scripts. Some of these scripts have a dependency to [jq](https://stedolan.github.io/jq/) and some others to [httpie](https://httpie.org/). 

## Build

The following snippets assume that you've opened a terminal and cd'ed to the the `host_lambda` directory.

###Build all Lambdas

```
mage -d ./mage build
```

which builds up to `max_concur_builds` lambdas concurrently. 
The value of `max_concur_builds` be set in the `./mage/conf.ini` file.  

###Build a Single Lambda

```
mage -d ./mage build <LambdaResourceName>
```

where LambdaResourceName is the resource name in the `template.yaml` file (e.g. `CheckCreationFunction`)

##Run Locally

TODO

##Debug Locally

Need build with env var `TN_DEBUG=true`. For all lambdas:

```
TN_DEBUG=true mage -d ./mage build
```

For a single lambda:

```
TN_DEBUG=true mage -d ./mage build <LambdaResourceName>
```

