# DB Migrate
This is command line utility to manage database migrations. This is developed for specific need where we have manage seprate migrations for diffrent schemas.

##Basic Rules for writing migration files
- Folers in which sql files exsists is treated as schema name. If the schema does not exist it is created so  the user should have permission to create the schema and also table.
- Only one level deep folders are checked.
- Every migration should have up and down scripts.
- The file name of the migration should be like
  < VersionNumber >-< Purpose >.{up/down}.sql
  Version Number should be integer, Purpose is a String without "-" or white spaces. It is recommended not to use any special chararecters in  the file names.
  Sample: 05-AddUser.up.sql,05-AddUser.down.sql

