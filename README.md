# Disk usage analyzer
A module and command line utility for getting info about disk usage.

# Usage instructions
Build a main application's entrance point:
```
go build danalyz.go
```

Run danalyz binary. In some cases it has to be executed with root privileges.

Required flags:
- __path__ - Path to directory to work with.
- __big-dirs__ - Big directories mode. Returns all subdirectories of given directory,
 which size is bigger than given max size (max-size flag). Either this flag or __dir_struct__
 flag has to be present.
- __dir-struct__ - File structure mode. Returns hierarchical representation of given directory
with given depth (depth flag), as well internal files sizes. Either this flag or __big_dirs__
flag has to be present.

Optional flags:
- __depth__ - Depth of walk for File structure mode. (default 3)
- __max-size__ - Files with size bigger than this value are considered as large files. (default 1000)

# Examples
Get subdirectories of /path/to/my/folder bigger than 500mb:
```
danalyz -path="/path/to/my/folder" -big-dirs -max-size=500
```

Get hierarchical structure of all subdirectories of /path/to/my/folder with max depth of 10:
```
danalyz -path="/home/max" -dir-struct -depth=10
```
