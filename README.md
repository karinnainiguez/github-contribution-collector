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
        repos: { owner-name: repo-name, second-owner: repo }
    ```
4. AWS Application Role with access to AWSLambda
    * Within your AWS account, create an application role and attach appropriate policies (including access to AWSLambda)
5. GitHub Access Key
    * The tool uses the GitHub API, which has rate limits.  The tool requires the rate limit increase granted by GitHub with a GitHub Access Key.  
    * Login to GitHub, and create the key.  You will access to that key when using the tool.


### CLI Usage

### Set Up Recurring Lambda Function
The attached zip file can be used as a Lambda function.  The lambda function will send an email to a provided (verified) email address.  Since the data includes all the contributions from the month until yesterday, it's recommended that you trigger the function once a month on the first of the month.  The instructions below include those paramaters, but feel free to trigger the function more often if you'd like to have the contributions so far for the month at a different date. 
1. Create a Lambda Function
2. Trigger Lambda Function Using 

### Monthly Tracking
Since the tool does not store any data for you, it may be important to see the amount of contributions per month for a given period of time.  If that's important to you, an easy solution is to track these in a spreadsheet with visualization built in.  Taking these extra steps will also ensure that you track contributions even after a user leaves the team and has been removed from the yaml file containing everyone's GitHub handles.

1. Use an excel template like the one provided. 
2. Once a month, when you receive the contributions for a given month, add those contributions in the data tab of the spreadsheet. 
3. In the reporting tab, right click the pivot table and select "refresh" to ensure that the data is updated on the graph. 
4. Save locally or in a shared drive depending on needs.

## Project Learning Goals

## Expectations

## Possible Future Enhancements
This tool was created during a short term internship.  However, additional functionality could be added depending on the time and resources available.  Below are some ideas for expansion. 
1. Integrate new GitHub API endpoints
2. Refactor for time complexity
3. Store data collected monthly in a shared location

