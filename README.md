# EKS OSS Contributions


## Usage

### Prerequisites
In order to use the tool as either a Command Line Interface locally, or a recurring AWS Lambda Function, there are several prerequisites the user will need in order to set up functionality. 
You will need: 

1. AWS Credentials
    * You'll need access to your account credentials.  _**DO NOT SHARE YOUR CREDENTIALS.  Instead, use environment variables where indicated.**_
    * Your accoung needs to be able to create an Application Role with access to Lambda, an S3 Bucket with a read-only object, and a Simple Email Service profile.
2. Simple Email Service Information - 
    * Create an AWS SES profile using the AWS Console. 
        * You will need access to your Username, Password, and Server Name.
    * Verify an email address from that profile.
        * Open the console to the Simple Email Service section.  
        * From the side menu on the left, select "Email Addresses"
        * Select "Verify a New Email Address" and enter the email address. 
        * Once the email is verified, select it and send a test email. 
3. S3 Bucket with Object (OPTIONAL - if you would rather not use an S3 bucket, you may still use the CLI tool with a local file, but will need to specify that with a flag)
    * From the AWS Console, select the S3 dashboard. 
    * Create a new bucket (It is recommended that you do NOT grant public read access to the bucket.  Instead, ensure that your specific User ID has read and write access)
    * Once the bucket is created, upload a file in yaml format that includes information about the following: 
        * 
4. AWS Application Role with access to AWSLambda
    * Within your AWS account, create an application role and attach appropriate policies (including access to AWSLambda)
5. GitHub Access Key
    * The tool uses the GitHub API, which has rate limits.  The tool requires the rate limit increase granted by GitHub with a GitHub Access Key.  
    * Login to GitHub, and create the key.  You will access to that key when using the tool.


### CLI Usage
### Set Up Recurring Lambda Function

## Project Planning

## Project Learning Goals

## Expectations

## Possible Future Enhancements

