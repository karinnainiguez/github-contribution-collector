# GitHub Contribution Collector


## Usage
This tool can be used to measure the contributions into open source projects of an entire team.  Originally created to facilitate metrics for the Amazon Kubernetes OSS team, it can also be applied to any other team by using the steps below.  To get started clone and install this repository inside of your GOPATH. (or use the ```go get github.com/karinnainiguez/github-contribution-collector``` command)

### Prerequisites
In order to use the tool as either a Command Line Interface locally, or a recurring AWS Lambda Function, there are several prerequisites the user will need in order to set up functionality. 
You will need: 

1. AWS Credentials
    * You'll need access to your account credentials.  _**DO NOT SHARE YOUR CREDENTIALS.**_
    * Your account needs to be able to create an Application Role with access to Lambda, an S3 Bucket with a read-only object, and a Simple Email Service profile.
2. Simple Email Service Information - 
    * Create an AWS SES profile using the AWS Console. 
        * You will need access to your Username, Password, and Server Name.
    * Verify an email address from that profile.
        * Open the console to the Simple Email Service section.  
        * From the side menu on the left, select "Email Addresses"
        * Select "Verify a New Email Address" and enter the email address. 
        * Once the email is verified, select it and send a test email. 
3. S3 Bucket with Object _(OPTIONAL - if you would rather not use an S3 bucket, you may still use the CLI tool with a local file, but will need to specify that with a flag)_
    * From the AWS Console, select the S3 dashboard. 
    * Create a new bucket (It is recommended that you do NOT grant public read access to the bucket.  Instead, ensure that your specific User ID has read and write access)
    * Once the bucket is created, upload a file in yaml format that includes information about the following: 
        * handles: a section with a list of GitHub handles of each member on your team.
        * orgs: a list of organizations where contributions should be tracked.  The tool will track contributions in all repositories within these organizations. 
        * repos: a list of repositories where contributions should be tracked.  Only include additional repositories that are not in organizations already tracked in the previous section.
    
    Below is an example of a yaml file with the necessary information: 
    ```yaml
        handles:
        - userName1
        - userName2
        - userName3
        orgs: 
        - orgname1
        - org-name2
        - third-org-name
        repos: { 
            repo1: owner1,
            repo2: owner2, 
            repo3: owner3,
        }
    ```
4. AWS Application Role with access to AWSLambda
    * Within your AWS account, create an application role and attach appropriate policies (including access to AWSLambda)
        * Example Trust Policy: 
        ```json
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "lambda.amazonaws.com"
                        },
                    "Action": "sts:AssumeRole"
                }
            ]
        }
        ```
5. GitHub Access Key
    * The tool uses the GitHub API, which has rate limits.  The tool requires the rate limit increase granted by GitHub with a GitHub Access Key.  
    * Login to GitHub, and create the key.  You will access to that key when using the tool.


### CLI Usage

Once you have cloned and installed the tool (and dependencies) on your machine, you will need to take the following steps in order to use the tool: 
1. Ensure that the yaml file containing your team's information is set up either in an S3 bucket, or locally (see below for appropriate flag).  The file must contain a section for handles, orgs, and repos as shown above. 
2. AWS credentials must be made accessible on your local machine.  This can be done through environment variables, or ```~/.aws/credentials``` file.  But the credentials must be for an account with access to the S3 bucket and SES account created in the prerequisites section. 
3. Environment Variables: the following environment variables must be set up in order for the tool to run properly: 
    * **GITHUBKEY**: Containing the API Key obtained through GitHub during the prerequisites section.
    * **SESVerifiedEmail**: An email address that has been verified.  This will be used as the **sender** and should be from the verified email address list on your AWS Simple Email Service profile. 
    * **S3BucketName**: OPTIONAL - the bucket name given to the resource created in prerequisites.  (in order to exclude this variable, you must use the ```--local-file``` flag during command)
    * **S3ObjectName**: OPTIONAL - the object name within the bucket. (in order to exclude this variable, you must use the ```--local-file``` flag during command)

Once those steps are complete, you'll be able to run commands with the tool in order to display a table in your terminal.  The table will contain contributions for your team.  Below are some supported commands and flags: 

To list all contributions for your team in the provided repositories so far during the current month: 
```
github-contribution-collector report
```
If you'd like to specify a ```--from``` and/or ```--until``` flag, the tool will display the contributions between those dates:
```
github-contribution-collector report --from 03-01-2018 --until 03-15-2018
```

Available optional flag ```--email``` will send an email containing the data to the specified (verified) email address:
```
github-contribution-collector report --email adrs@example.com
```

To use a local yaml file with team imformation instead of an S3 bucket with an object, use the optional ```--local-file``` flag, and provide the complete path: 
```
github-contribution-collector report --local-file /Users/username/folder/filename.yaml
```
 

### Set Up Recurring Lambda Function
The attached zip file can be used as a Golang Lambda function.  The lambda function will send an email to a provided (verified) email address.  Since the data includes all the contributions from the month until yesterday, it's recommended that you trigger the function once a month on the first of the month.  The instructions below include those parameters, but feel free to trigger the function more often if you'd like to have the contributions so far for the month at a different date. 
1. Create a Lambda Function.  
    * From the AWS Lambda console, select "Create Function" button.  
    * Name the function, select the Go 1.x runtime, and use the existing role created during the prerequisites (that will ensure that your lambda function has access to all other necessary AWS resources). 
    * In the Function code section, upload the zip file in the repository and name the Handler "github-contribution-collector"
    * In the Environment Variable section, ensure that the following Environment Varibles are defined: 
        * **GITHUBKEY**: Containing the API Key obtained through GitHub during the prerequisites section.
        * **SESVerifiedEmail**: An email address that has been verified.  This will be used as the **sender** and should be from the verified email address list on your AWS Simple Email Service profile. 
        * **S3BucketName**: The bucket name given to the resource created in prerequisites.  (this is not optional in this process)
        * **S3ObjectName**: The object name within the bucket. (this is not optional in this process)
    * In the Basic settings section, increase the timeout to 15 minutes.
2. Trigger Lambda Function Using CloudWatch Event. 
    * On the AWS Lambda Console, there is an option to add triggers from a list.  Click on CloudWatch Events.
    * Select "Create a new rule" and enter a name and description. 
    * For the rule type, select schedule expression and provide the following: 
        * cron(0 6 1 1/1 ? *)
    * Click add. Save changes to the Lambda function. 
    * Configure a new test event from template Amazon CloudWatch. 
    * Test the function.

### Monthly Tracking
Since the tool does not store any data for you, it may be important to see the amount of contributions per month for a given period of time.  If that's important to you, an easy solution is to track these in a spreadsheet with visualization built in.  Taking these extra steps will also ensure that you track contributions even after a user leaves the team and has been removed from the yaml file containing everyone's GitHub handles.

1. Use an Excel template like the one provided in this repository.
2. Once a month, when you receive the contributions for a given month, add those contributions in the data tab of the spreadsheet. 
3. In the reporting tab, right click the pivot table and select "refresh" to ensure that the data is updated on the graph. 
4. Save locally or in a shared drive depending on needs.

## Possible Errors
During the use of the CLI tool, the user may get one of the following errors: 

Rate API Exceeded: 
```
[✖] 403 API rate limit of 5000 still exceeded until ... not making remote request. [rate reset in 51m16s]
```
In the case above, the user has made too many API calls.  Rate limits are set by GitHub and are tied to a user.  The two possible fixes include waiting for 1 hr until the rate limit is reset to run the program again, or using another team member's GitHub API key.  This program does not support rotating API Keys internally. 

Unsupported Date: 
```
[✖]  parsing time notadate: month out of range
```
the user has an option of providing a from and until flag argument with the date from and until that contributions should be tracked.  The supported format is "MM-DD-YYYY" but the user also has the option to exclude those flags and get all contributions so far in the current month.

Server Error: 
```
[✖]  GET: 502 Server Error []
```
The server is unavailable.  The local machine must be connected to the internet.  Ensure connection and try again. 

## Expectations
The end user was looking for a way to receive an email notification once a month, with the information regarding a team's contributions into Open Source Software.  Original functionality was to measure contributions including issues raised and pull requests submitted, from one repository.  During implementation and updates with the user, the functionality was extended to include **all** repositories under multiple organizations, as well as individual repositories specified.

## Possible Future Enhancements
This tool was created during a short term internship.  However, additional functionality could be added depending on the time and resources available.  Below are some ideas for expansion. 
1. **Integrate new GitHub API endpoints.** When this project was developed, the limitation of checking a user's contributions per repository made the functionality very expensive.  Since GitHub has rate limits for API calls, it would be ideal to have the program collect contributions in an entire organization with one call per user.  Future enhancements to the API, or the GraphQL API may make that call possible, significantly reducing the cost of adding a new team member or a new organization to our list of contributions.  This request has been brought up in the API forum, and may be available soon.
2. **Deliver an attachment with the monthly email.**  Since the email contains an html table with the contributions, it makes it easy for the user to forward that email to interested parties.  However, if the user wanted to manipulate the data or visualize progress within the team, it may be easier to do that with an attached csv file containing the contributions. 
3. **Store data collected monthly in a shared location:** Currently, the user would have to store information locally or in a shared drive in order to see the progress of the team's contributions over time.  If this data were to be stored in a database, or a cloud resource, it would facilitate functionality to display that progress without the cost of making all of those API calls every time, or tracking those monthly contributions locally. 
4. **Test Coverage.**  Testing the tool will be very important when making enhancements.  However, most of the functionality revolves around external resources and APIs.  Because rate limits would be greatly impacted, real API calls should not be sent during tests, but instead mocked for greater efficiency. 

