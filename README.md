# DB Migrate
This is command line utility to manage database migrations. This is developed for specific need where we have manage seprate migrations for diffrent schemas.

## Basic Rules for writing migration files
- Folders in which sql files exsists is treated as schema name. If the schema does not exist it is created so  the user should have permission to create the schema and also table.
- Only one level deep folders are checked.
- Every migration should have up and down scripts.
- The file name of the migration should be like
  < VersionNumber >-< Purpose >.{up/down}.sql
  Version Number should be integer, Purpose is a String without "-" or white spaces. It is recommended not to use any special chararecters in  the file names.
  Sample: 05-AddUser.up.sql,05-AddUser.down.sql
  
## To Test
- checkout this project
- run ```go install``` in the project.
- Keep the Database URL handy.
- create a directory with names as specified above.(A sample folder is found in this repo)
- Go the sample folder
- Run the following command to Upgrade the database.
```
dbmigrate up --databaseURL postgres://postgres:postgres@localhost:5432/test --directory .
```


## To Release 
You need to have [goreleaser](https://goreleaser.com/) installed on your machine
export GITHUB_TOKEN=`YOUR_GH_TOKEN`
git tag -a v0.1.3 -m "Vulnerability Fixes"
git push origin v0.1.3
goreleaser
