
# What is iis-hero? 

iis-hero is a CLI management tool for IIS servers. It provides access and management capabilities to remote IIS servers.
It supports a wide range of operations for application pools, web sites, applications, and virtual directories.
In addition to IIS, it also provides support for Windows services and file operations. You can manage Windows services and perform operations such as file transfer between remote servers.


## Getting started

Assuming that Go is installed, if it is not

[Download](https://go.dev/doc/install) and install Go.

```
go install github.com/nub06/iis-hero@latest
```

Check your `$GOPATH/bin` 

You can rename the binary name if you don't want to use CLI as `iis-hero`


# Feautures

#### Available Commands 

- [execute](#execute) (Execute custom powershell commands on target computer)
- [start](#backing-up-iis-configuration) (Start all Internet services on target computer.)
- [stop](#backing-up-iis-configuration) (Stop all Internet services on target computer)
- [reset](#backing-up-iis-configuration) (Stop and then restart all Internet services on the target computer)
- [config](#config-commands-usage) (Management IIS configuration)
  - [backup](#backing-up-iis-configuration) (Backup IIS configuration)
    - [remove](#removing-iis-configuration-backup-file) (Remove IIS configuration backup file)
  - [clear](#clearing-all-iis-configuration) (Clear all IIS configuration data.)
  - [restore](#restoring-iis-configuration-from-backup) (Restore IIS configuration from backup.)
  - [list](#list-the-configuration-backups) (List IIS configuration backups.) 
- [pool](#iis-application-pool-commands-usage) (Manage Application Pools on IIS)
  - [change](#changing-iis-applicationpool-properties) (Change the properties of an Application Pool on IIS)
  - [create](#creating-iis-applicationpool)  (Create an Application Pool on IIS)
  - [list](#listing-iis-application-pools-and-properties) (List Application Pools and properties on IIS)
  - [remove](#deleting-an-iis-application-pool) (Remove an Application Pool on IIS)
  - [start](#starting-an-iis-application-pool) (Start the Application Pool on IIS)
  - [stop](#stopping-an-iis-application-pool) (Stop the Application Pool on IIS)
  - [restart](#restarting-an-iis-application-pool) (Restart the Application Pool on IIS)
  - [state](#current-state-of-the-iis-application-pool) (Display the status of an Application Pool on IIS)
- [site](#iis-site-commands-usage) (Manage Web Sites on IIS)
  - [change](#changing-iis-website-properties) (Change the properties of Web Site on IIS)
  - [create](#creating-an-iis-website) ( Create a Web Site on IIS)
  - [list](#listing-iis-websites-and-properties) (List Web Sites and properties on IIS)
  - [remove](#removing-an-iis-website)  (Remove a Web Site on IIS)
  - [start](#starting-an-iis-website) ( Start the Web Site on IIS)
  - [stop](#stopping-an-iis-website) (Stop the Web Site on IIS)
  - [state](#iis-site-state) ( Display the status of a Web Site on IIS)
- [vdir](#iis-virtual-directory-commands) (Manage virtual directories on IIS)
  - [change](#changing-iis-virtual-directory-settings) (Change the properties of Virtual Directory on IIS)
  - [create](#creating-iis-virtual-directory) (Create a Virtual Directory on IIS)
  - [list](#listing-iis-virtual-directories-and-properties) (List Virtual Directories and properties on IIS)
  - [remove](#removing-iis-virtual-directory) (Remove a Virtual Directory on IIS)
- [app](#iis-application-commands) (Manage IIS Applications)
  - [change](#changing-iis-application) (Change the properties of an Application on IIS)
  - [create](#creating-iis-application) (Create an Application on IIS)
  - [list](#listing-iis-applications) (List Applications and properties on IIS)
  - [remove](#removing-iis-application) (Remove an Application on IIS)
- [folder](#folder-commands) (Manage folders on target computer)
  - [backup](#backup-command) (Create a backup of the specified folder on the target computer)   
  - [list](#list-command) (List folders on given path)
  - [push](#push-command) (Perform file copying from your local computer to the target computer) 
  - [pull](#pull-command) (Perform file copying from target computer to your local computer)  
- [winsvc](#windows-service-commands) (Manage Windows Services)     
  - [change](#changing-windows-service-properties) (Change the properties of the existing Windows Service)
  - [create](#creating-windows-service)(Create a new Windows Service)
  - [list](#listing-windows-service-and-properties) (List Windows Service and properties)
  - [remove](#deleting-windows-service-properties) (Remove a Windows Service with given name)
  - [start](#starting-windows-service) (Start the Windows Service)
  - [stop](#stopping-windows-service) (Stop the Windows Service)
  - [restart](#restarting-windows-service) (Restart the Windows Service)
  - [state](#current-state-of-the-windows-service) (Display the status of a Windows Service)
- [login](#to-use-the-this-cli-first-you-need-to-specify-the-target-computer)   
  - [clear](#to-use-the-this-cli-first-you-need-to-specify-the-target-computer) (Clears the current login credentials)   
  - [cred](#to-use-the-this-cli-first-you-need-to-specify-the-target-computer)  (Display current credentials) 



# Usage 

## To use the this CLI, first you need to specify the target computer.


- If the target computer is already your own computer, you can use `local` or `localhost` on this command

```
iis-hero login -c local
iis-hero login -c localhost
```

- If you have access to the target computer with the same credentials, use this command. It will connect directly with the ComputerName.

```
iis-hero login -c ComputerName
```

- If you have access to the target computer with different credentials and WinRM configuration is properly set up, use this command. It will connect to the target computer using the credentials you specify.

```
iis-hero login -c ComputerName -d Domain -u UserName -p Password
```

- Once you have specified the target computer, it will remain the same until you change it. You don't need to specify it again every time, even if you restart your computer.

- Your login credentials will be stored encrypted in the `%appdata%\iis-hero`. If you don't want to store them, you can use the `iis-hero login clear` command in the program every time you finish your work.

- You can view the credentials used by the application by using the `iis-hero login cred` command.

- Use `iis-hero login cred -f` if you want to view your password without asterisk.


# Config Commands usage

## Backing up IIS configuration.

- Usage:
  - `iis-hero config backup [flags]`
  - `iis-hero config backup <foldername> [flags]`
  - `iis-hero config backup [command]`

- Available Commands:
  - remove   Remove IIS configuration backup

- Flags:
  - -h, --help   help for backup

- This command creates a backup folder in the `C:\Windows\System32\inetsrv\backup`
- This command can be used with or without an argument. If you specify a folder name as an argument, it will create a folder with the specified name. 

- If you don't use any argument in this way, it will create a backup folder in the following format.  `IISConfigBackup_yyyyMMdd_HHmm` 

- Usage example:

```
iis-hero config backup
iis-hero config backup LastBackup
```

### Removing IIS configuration backup file.

- You can provide a folder name as an argument to this command. It deletes the backup file with the name of the folder you entered as an argument. If you want to delete all backups, use the -a, --all flag instead of passing an argument.

- Usage:
  - `iis-hero config backup remove <foldername> [flags]`

- Aliases:
  - remove, rm

- Flags:
  - -a, --all    Flag for the remove all backup files
  - -h, --help   help for remove

- Usage example:  
```
iis-hero config backup rm LastBackup
iis-hero config backup remove LastBackup
iis-hero config backup rm -a
iis-hero config backup rm --all
```

## Clearing all IIS configuration.

- This command clears all IIS configuration.

- Usage:
  - `iis-hero config clear [flags]`

- Flags:
  -h, --help   help for clear

- Usage example:

```
iis-hero config clear
```

## Restoring IIS configuration from backup.

- Usage:
  - `iis-hero config restore [flags]`

- Flags:
  - -h, --help      help for restore
  - --latest        Restores the most recently created backup file
  - --name string   Restores the backup file with the specified name

- Usage example:  

```
iis-hero config restore --latest
iis-hero config restore --name IISConfigBackup_20230425_1148
```

- If you use this command with the `--latest` flag, it restores the latest backup folder. However, with the `--name flag`, you can specify a specific version that you want to restore.


## List the configuration backups.

- This command lists the backup folders that have been previously taken in the `C:\Windows\System32\inetsrv\backup` 

- Usage:
  - `iis-hero config list [flags]`

- Aliases:
  - list, ls

- Flags:
  - -h, --help   help for list

- Usage example:
```
iis-hero config ls 
iis-hero config list
```

# IIS Site Commands usage
## Changing IIS Website Properties.

- Usage:
  - `iis-hero site change <sitename> [flags]`

- Flags:
     - --bind string      Set new Binding Value for the Web Site
     - -h,  --help             help for change
     -  --path string      Set new Physical Path Value for the Web Site
     -  --pool string      Set new Application Pool Name for the Web Site
     -  --preload string   Set new PreLoadEnabled Value for the Web Site
     -  --rename string    Rename Web Site

- You can use all the flags together if you want, or you can only use the flag of the value you want to change.

- Usage examples:
```
iis-hero site change MySite --path D:\Applications\MySite
iis-hero site change MySite --path D:\Applications\MySite  --pool MyAppPool
iis-hero site change MySite --path D:\Applications\MySite  --pool MyAppPool --preload false
iis-hero site change MySite --path D:\Applications\MySite  --pool MyAppPool --preload false --bind siteaddress.test
iis-hero site change MySite --pool MyAppPool 
```
## Creating an IIS Website.

- Usage:
  - `iis-hero site create <sitename> [flags]`

- Flags:
     -  --bind string      Identify Binding Path Value for the Web Site
     -  -h, --help         help for create
     -  --path string      Identify Physical Path Value for the Web Site
     -  --pool string      Identify Application Pool Name for the Web Site
     -  --preload string   Identify PreLoadEnabled Value for the Web Site (default "true")


 - Usage examples:
 ```
 iis-hero site create NewSite --pool DefaultAppPool --bind testaddress.test --path D:\Application\NewSite 
 iis-hero site create NewSite --pool DefaultAppPool --bind testaddress.test --path D:\Application\NewSite --preload false
```
## Listing IIS Websites and properties

- Usage:
  - `iis-hero site list [flags]`

- Aliases:
  - list, ls

- Flags:
  - -a, --all       -a , --all returns all sites
  - -h, --help      help for list
  - -s, --stopped   -s, --stopped returns stopped sites
  - -n, --name string   Specify the name of the website you want to list

- If you want to list a single web site. Use  `-n` flag.
- If you want to list only the running websites, don't use any flag.
- If you want to list only the stopped websites, use the `-s` flag. 
- If you want to list all websites, use the `-a` flag.

 - Usage examples:

```
iis-hero site list
iis-hero site ls -a
iis-hero site ls
iis-hero site list -s
iis-hero site ls -s
iis-hero site list -a
iis-hero site ls -n MySite
iis-hero site ls --name MySite

```
## Removing an IIS Website.

- Usage:
  - `iis-hero site remove <sitename>`

- Aliases:
  - remove, rm

- Flags:
  - -h, --help   help for remove

- Usage examples:

```
iis-hero site remove NewSite 
iis-hero site rm NewSite 
```
## Starting an IIS Website.

- Usage:
  - `iis-hero site start <sitename>`

- Flags:
  - `-h, --help   help for start`

- Usage examples:
```
iis-hero site start MySite
```
## Stopping an IIS Website.

- Usage:
  - `iis-hero site stop <sitename>`

- Flags:
  - -h, --help   help for stop

- Usage examples:
```
iis-hero site stop MySite
```
## IIS Site State

- Usage:
  - `iis-hero site state <sitename>`

- Aliases:
  - state, st

- Flags:
  - -h, --help   help for state

- Usage examples:
```
iis-hero site state MySite
iis-hero site st MySite
```
# IIS Application Pool Commands
## Changing IIS ApplicationPool Properties.
- Usage:
  - `iis-hero pool change <poolname> [flags]`
  
- Flags:
     -  --autostart string    Set Autostart for application pool
     - -h, --help             help for change
     -  --idleaction string   Set Idle-Timeout Action for application pool
     -  --idleminute string   Set Idle-Timeout(minutes) for application pool
     -  --pipeline string     Set Managed Pipeline Mode for application pool
     -  --rename string       Set new name for application pool
     -  --clr string      Set Runtime version for application pool
     -  --startmode string    Set Start Mode for application pool 

- You can use all the flags together if you want, or you can only use the flag of the value you want to change.

- .NETCLR Version ( Runtime version)

 - Use `--clr v2` for `.NETCLR Version v2.0`
 - Use `--clr v4` for `.NETCLR Version v4.0` 
 - Use `--clr 0` for `No Managed Code`

- Usage examples:
```
iis-hero pool change MyAppPool --clr 0
iis-hero pool change MyAppPool --idleminute 0
iis-hero pool change MyAppPool --clr v4 --idleminute 0 --pipeline Classic 
iis-hero pool change MyAppPool --clr v2 --idleminute 5 --pipeline Integrated --idleaction Suspend
iis-hero pool change MyAppPool --clr 0 --idleminute 20 --pipeline Integrated --idleaction Terminate --autostart false
iis-hero pool change MyAppPool --clr 0 --idleminute 20 --pipeline Integrated --idleaction Terminate --autostart false --startMode OnDemand
```
## Creating IIS ApplicationPool.

- Usage:
  - `iis-hero pool create <poolname> [flags]`
- Flags:
    -  --autostart string    Set Autostart for application pool (default "true")
    -  -h, --help            help for create
    -  --idleaction string   Set Idle-Timeout Action for application pool (default "Suspend")
    -  --idleminute string   Set Idle-Timeout(minutes) for application pool (default "0")
    -  --pipeline string     Set Managed Pipeline Mode for application pool (default "Integrated")
    -  --clr string          Set Runtime version for application pool e.g: --clr v4, --runtime v2 (default: "No Managed Code")
    -  --startmode string    Set Start Mode for application pool (default "AlwaysRunning")

- .NETCLR Version ( Runtime version)
 - Use `--clr v2` for `.NETCLR Version v2.0`
 - Use `--clr v4` for `.NETCLR Version v4.0` 
 - Use `--clr 0` for `No Managed Code`


- Usage examples:
```
iis-hero pool create NewAppPool
iis-hero pool create NewAppPool --startmode OnDemand
iis-hero pool create NewAppPool --startmode OnDemand --clr v4
iis-hero pool create NewAppPool --clr v2 --idleaction 5

```
## Listing IIS Application Pools and properties.

- Usage: 
  - `iis-hero pool list [flags]`
- Aliases:
  - list, ls
- Flags:
  - -a, --all       List all application pools
  - -h, --help      help for list
  -  -s, --stopped  List only stopped application pools

- If you want to list only the running application pools, don't use any flag.
- If you want to list only the stopped application pools, use the `-s` flag. 
- If you want to list all application pools, use the `-a` flag. 

 - Usage examples:
```
iis-hero pool list
iis-hero pool list -s
iis-hero pool list -a
iis-hero pool ls -a
iis-hero pool ls -s
iis-hero pool ls
```

## Removing an IIS Application Pool.

- Usage:
  - `iis-hero pool remove <poolname> [flags]`
- Aliases:
  - remove, rm
- Flags:
  - -f, --force   Force flag to delete an Application Pool
  - -h, --help    help for remove 

- If there is a Web Site belonging to the Application Pool that you want to delete, but you still want to delete you should use the `--force` flag to delete that Application Pool.

- Usage examples:

```
iis pool remove NewAppPool
iis pool remove NewAppPool -f
```

## Starting an IIS Application Pool.

- Usage:
  - `iis-hero pool start <poolname>`
- Flags:
  - -h, --help   help for start  
- Usage examples:
```
iis-hero pool start MyAppPool
```
## Stopping an IIS Application Pool.

- Usage:
  - `iis-hero pool stop <poolname>`

- Flags:
  - -h, --help   help for stop

- Usage examples:
```
iis-hero pool stop MyAppPool
```
## Restarting an IIS Application Pool.


- Usage:
  - `iis-hero pool restart <poolname> `

- Aliases:
  - restart, reset

- Flags:
  - -h, --help   help for restart
  - -a, --all    Restart all Application Pools

- Usage example: 
```
iis-hero pool restart MyAppPool
iis-hero pool restart -a
iis-hero pool reset -a
```
## Current state of the IIS Application Pool


- Usage:
  - `iis-hero pool state <poolname>`

- Aliases:
  - state, st
- Flags:
  - -h, --help   help for state
- Usage example:
```
iis-hero pool state MyAppPool
iis-hero pool st MyAppPool
```

# IIS Virtual Directory Commands.
## Changing IIS Virtual Directory Settings.
- Usage:
  - `iis-hero vdir change <name> [flags]`

- Flags: 
  - --path string           Set a new Physical Path Value for the Virtual Directory
  - --rename string         Set a new Name  Value for the Virtual Directory

- You can use all the flags together if you want, or you can only use the flag of the value you want to change.
- Usage examples: 

```
iis-hero vdir change MyVirtualDirectory --path D:\Applications\MyVDir
iis-hero vdir change MyVirtualDirectory --rename MyNewVDir
```

## Creating IIS Virtual Directory.

- Usage: 
  - `iis-hero vdir create <name> [flags]`
- Flags:
  - -h, --help          help for create
  - --path string   Identify Physical Path Value for the Virtual Directory
  - --site string   Identify Web Site Name Value for the Virtual Directory
- Usage examples: 
```
iis-hero vdir create NewVirtualDirectory --site MyWebSite --path D:\Application\NewVirtualDir 
```
## Listing IIS Virtual Directories and properties.

- Usage:
  - `iis-hero vdir list [flags]`

- Flags:
  - -h, --help   help for list
- Usage example:
```
iis-hero vdir list
```

## Removing IIS Virtual Directory.
- Usage:
  -` iis-hero vdir remove <name> [flags]`
- Aliases:
  - remove, rm  
- Flags:
  - -h, --help          help for remove
  - --site string   Specify the Web Site Name of the Virtual Directory
  - --app string    Specify the Application Name of the Virtual Directory

- 
If you don't specify the Site and Application, the iis-hero will automatically search and delete the Virtual Directory you specified with `<name>`. However, if there are multiple virtual directories with the same name on the IIS server, the cli will prompt you to specify the Site and Application name of the Virtual Directory you want to delete using the `--site` and `--app` flags.

- Usage examples:
 ```
iis-hero vdir remove NewVirtualDirectory --site MyWebSite
iis-hero vdir remove NewVirtualDirectory 
iis-hero vdir rm NewVirtualDirectory --site NewWebsite
iis-hero vdir rm NewVirtualDirectory --site NewWebsite --app MyApp
``` 

# IIS Application Commands.

- If you do not specify a site with the `--site` flag, the CLI will search for the site of the application you specified as `<application name>`. If it cannot find it, it will show you a warning. In that case, use the `--site` flag to specify the site name.

## Changing IIS Application.

- Usage:
  - `iis-hero app change <application name> [flags]`

- Flags:
  - -h, --help             help for change
  - --path string      Set a new Physical Path for Application
  - --pool string      Set a new Application Pool for Application
  - --preload string   Set a PreloadEnabled value for Application
  - --rename string    Set a new name for Application
  - --site string      Identify Site Name of the Application


- Usage example:  

```
 iis-hero app change MyApp --pool NewAppPool 
 iis-hero app change MyApp --pool NewAppPool2 
 iis-hero app change MyApp --rename MyApplication --site NewWebSite
 ```
## Creating IIS Application.

- Usage:
  - `iis-hero app create <application name> [flags]`

- Flags:
  - -h, --help      help for create
  - --path string   Identify Physical Path for Application
  - --pool string   Identify Application Pool for Application
  - --site string   Identify Site Name for Application

- Usage example:

```
 iis-hero app create NewApp --pool NewAppPool --site NewWebSite --path D:\NewApp
 ```

## Removing IIS Application

- If you want to delete an application, use the `iis-hero app remove` command. If this application contains other applications, the application will give you an error message stating that it has sub-applications. If you still want to delete it, use the force flag. `iis-hero app remove -f` If you use this flag, the specified application and all applications contained within it will be deleted."

- Usage:
  - `iis-hero app remove <application name> [flags]`

- Aliases:
  - remove, rm  

- Flags:
  - -f, --force         Force flag for the Remove Application
  - -h, --help          help for remove
  - --site string       Identify Site Name for Application

- Usage example: 
 ```
iis-hero app remove MyApp -f
iis-hero app rm MyApp --force
iis-hero app rm MyApp
  ```

## Listing IIS Applications.

- Usage:
  - `iis-hero app list [flags]`

- Flags:
  - -h, --help   help for list
- Usage example:
 ```
iis-hero app ls 
iis-hero app list
```

 # Windows Service Commands.
 - When specifying the `<servicename>` parameter in all commands except for the `create` command, you can use either the service name or the display name as the value.
## Changing Windows Service Properties.
- Usage :
  - `iis-hero winsvc change <servicename> [flags]`
- Usage examples:
 ```
iis-hero winsvc change MyWinService --description "New Description of NewWinService" --displayname "New Win Service" 
iis-hero winsvc change MyWinService --description "New Description of NewWinService" --startup Automatic
iis-hero winsvc change MyWinService --description "New Description of NewWinService" --startup Automatic --rename RenamedWinService
```
- Flags:
  - --description string   Set new Description value for Windows Service
  - --displayname string   Set new DisplayName value for Windows Service
  - --exepath string       Set new Executable Path value for Windows Service
  - -h, --help             help for change
  - --startup string       Set new StartupType value for Windows Service
## Creating Windows Service.

- The `<service name>` parameter corresponds to the service name value of the Windows service that will be created.

- Usage 
  - `iis-hero winsvc create <servicename> [flags]`

- Usage examples:
 ```
iis-hero winsvc create NewWinSvc  --displayname "New Win Svc" --exepath "D:\Application\MyWinSvc\mysvc.exe"
iis-hero winsvc create NewWinSvc  --displayname "New Win Svc" --description "The service is used for MyApp" --exepath "D:\NewWinSvc\mysvc.exe"
iis-hero winsvc create NewWinSvc  --displayname "New Win Svc" --description "The service is used for MyApp" --exepath "D:\NewWinSvc\mysvc.exe" --startup Automatic
iis-hero winsvc create MyWinService --exepath "D:\NewWinSvc\mysvc.exe" 
```
## Deleting Windows Service Properties.
- Usage:
  - `iis-hero winsvc remove <servicename> [flags]`
- Aliases:
  - remove, rm
- Flags:
  - -h, --help   help for remove
- Usage examples:
 ```
iis-hero winsvc remove MyWinService 
iis-hero winsvc rm "My Win Service" 
 ```
## Stopping Windows Service.
- Usage:
  - `iis-hero winsvc stop <servicename> [flags]`

- Flags:
  - -h, --help   help for stop

- Usage examples: 

```
iis-hero winsvc stop MyWinService

```
## Starting Windows Service.
- Usage:
  - `iis-hero winsvc start <servicename> [flags]`

- Flags:
  - -h, --help   help for start

- Usage examples: 

```
iis-hero winsvc start MyWinService

```

## Restarting Windows Service


- Usage:
  - `iis-hero winsvc restart <servicename> [flags]`

- Aliases:
  - restart, reset

- Flags:
  - -h, --help   help for restart

- Usage examples: 

```
iis-hero winsvc restart MyWinService

```

## Listing Windows Service and Properties.
- Usage:
  - `iis-hero winsvc list <servicename> [flags]`

- Aliases:
  - list, ls

- Flags:
  - -h, --help   help for list
 
- Usage examples: 
```
iis-hero winsvc list MyWinService
iis-hero winsvc ls MyWinService
```
## Current state of the Windows Service.

- Usage:
  - `iis-hero winsvc state <servicename> [flags]`

- Flags:
  - --h, --help   help for state
- Usage examples: 
```
iis-hero winsvc state MyWinService
```

# Folder Commands.

## Backup Command

- This command takes a backup of a specified folder on the target computer. If you use the `--dest` flag to specify a path, it creates a backup folder at the specified path. If you don't use this flag, it will create a backup at the default path of `D:\Backups.`

- Usage:
  - `iis-hero folder backup <source folder path> [flags]`

- Flags:
  - -h, --help          help for backup
  - -d, --dest string   Specify backup folder destination (default "D:\\Backups")

- Usage example:

```
iis-hero folder backup D:\MyApplication
iis-hero folder backup D:\MyApplication --dest D:\NewFolder
```

## List Command

- Usage:
  - `iis-hero folder list <folderpath> [flags]`

- Aliases:
  - list, ls

- Flags:
  - -h, --help   help for list

- Usage example:

```
iis-hero folder ls D:\MyApplication
iis-hero folder list D:\MyApplication
```

 ## Push Command

 - This command performs file copying from your local computer to the target computer. Use the `--local` flag to specify the location of the folder you want to copy on your local computer, and use the `--target` flag to specify the destination where you want to copy the file on the target computer.


- Usage:
  - `iis-hero folder push [flags]`
- Flags:
  - -h, --help            help for push
  - --local string    Specify file/folder path on local computer
  - --target string   Specify destination folder path on target computer

- Usage example: 

``` 
iis-hero folder push --local D:\Folder --target D:\DestFolder
iis-hero folder push --local D:\try.txt --target D:\DestFolder
```


## Pull Command

 - This command performs file copying from your target computer to the local computer. Use the `--target` flag to specify the location of the folder you want to copy on your target computer, and use the `--local` flag to specify the destination where you want to copy the file on the local computer.


- Usage:
  - `iis-hero folder pull [flags]`
- Flags:
  - -h, --help            help for push
  - --local string    Specify destination folder path on local computer
  - --target string   Specify file/folder path on target computer

- Usage example: 

``` 
iis-hero folder pull --local D:\LocalDestFolder  --target D:\DestFolder
iis-hero folder pull --local D:\LocalDestFolder  --target D:\hero.json
```