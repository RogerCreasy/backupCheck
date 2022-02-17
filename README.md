# backupStatusCheck
Go program that validates backup files on a server<br>
This program is meant to be compiled and installed as a daemon on a Linux server. It checks that a file exists, and that it is newer than 7 days. If either criteria fails, a message is sent to a Google Space

## Requirements

  * Google Workspace account
  * A Space with a webhook
  * A Linux server (tested on Fedora

### Go packages
    none - this program uses only standard library packages

### Generate your Space webhook

## Create the Server-Side configuration
You need to create 2 configuration files - one for the webhook URL, another for the files you wish to check.

### The config.json file
In your editor, create a file named config.json. There is a config.json example file you can use as a starting point. This file is located in the serverFiles directory. Copy the webhook described above. Paste the ID in the config file as the appropriate json value.

```
   {
       "webhook": "Paste the webhook URL here",
   }
```

Upload config.json to the /etc/backupCheck directory on your server. Create the backupCheck directory, if it doesn't exist. <br>

### The backupFiles.json file
Add absolute file paths for each file you wish to check within.<br><br>

Upload backupFiles.json to the /etc/backupCheck directory on your server. Create the backupCheck directory, if it doesn't exist. <br>

## Development
TODO - explain compiling, etc

## Setting up the service and schedule
TODO - explain
