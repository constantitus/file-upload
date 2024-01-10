*this is pretty much a work in progress*

## File Upload
Web application for uploading and managing files. \
Can upload multiple files, overwrite existing files if needed, download, rename and delete files from storage. \
Users are stored in a database.

## Usage
Clone the repo, cd into it and run `make`. The binary is located at *build/upload*. \
Run the app once to create the database. To manually add users, run `make dbcli` and then use the `build/dbcli` tool. \
You can provide the port as an argument when running the app.

## TODO
- [x] Tool to manage users in the database
- [ ] Admin panel to manage users
