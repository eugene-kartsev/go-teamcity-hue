##### What is 'go-teamcity-hue'?
It is an application written in [golang](https://golang.org/) to make it possible to connect a [TeamCity](www.jetbrains.com/teamcity) build status with a [HUE-Phillips](http://www.developers.meethue.com/) lamps' signals. So that when your build goes "red", your HUE lamp turns red too.


##### How do I make it work?
The following steps describe the installation process for MacOS, but they are pretty much similar for every system.

1. You're going to need a [TeamCity installation](https://confluence.jetbrains.com/display/TCD9/Installing+and+Configuring+the+TeamCity+Server).

2. Make sure you can log-in into TeamCity REST Api and see your latest build status. An example for the TeamCity Api url would be:

 ```
 https://<HOST>/httpAuth/app/rest/buildTypes/id:<BUILD_ID>/builds?locator=branch:<BRANCH>&count=1
 ```
   
   and here is an example of a TeamCity REST Api response:
   
   ```xml
  <builds count="1" href="/httpAuth/app/rest/buildTypes/id:ID/builds?locator=branch:DEV&count=1"
          nextHref="/httpAuth/app/rest/buildTypes/id:ID/builds?locator=count:1,start:1,branch:DEV">
          
    <build id="139992" buildTypeId="ID" number="2.1-beta.53"
           status="SUCCESS" state="finished" branchName="DEV"
           defaultBranch="true" href="/httpAuth/app/rest/builds/id:139992" 
           webUrl="https://HOST/viewLog.html?buildId=139992&buildTypeId=ID"/>
           
  </builds>
   ```
   
   If you need more info about how the TeamCity REST Api works, check [their documentation](https://confluence.jetbrains.com/display/TW/REST+API#RESTAPI-BuildLocator).
   
3. You're going to need a HUE lamp, similar to this one: [hue bloom](https://www.google.ca/search?q=hue+bloom). The lamp connects to your wi-fi network and listens to your "commands" via simple http-request-response web server (check their [REST api](http://www.developers.meethue.com/philips-hue-api) - it's pretty straightforward).

4. Build 'go-teamcity-hue' from the source code (using [golang compiler](https://golang.org/doc/code.html)) or download the latest binary release from this repository.

5. OK, now you are ready to start the application for the first time. Your target directory should look like this:

 ```sh
 user$ ls
   go-teamcity-hue
 ```

6. Start the application for the first time - the configuration file teamplate will be created:

 ```sh
 user$ ./go-teamcity-hue
 A new configuration file has been created.
 Modify it and restart the application.
 ```

7. Check that the configuration file exists:

 ```sh
 user$ ls
  config		go-teamcity-hue
 ```

8. The configuration file should look similar to this:

 ```sh
 user$ cat config
    {
        "version": "0.01",
        "hueNodes": [{
            "id": "hue1",
            "url": "<HUE_URL>"
        }],
        "teamcityNodes": [{
            "id": "tc1",
            "url": "<TEAMCITY_BUILD_URL>",
            "login": "<USER_LOGIN>",
            "password": "<USER_PASSWORD>",
            "interval": 10
        }],
        "map": [{
            "hueId": "hue1",
            "teamcityIds": ["tc1"]
        }]
    }
 ```

9. Modify the file accordingly to your settings.

10. After your configuration file is ready, start the application once again. You should see the messages in your console and your HUE lamp should start changing its color depending on TeamCity build status.
