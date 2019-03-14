# awparse


Parser for automated weather station data 
archive from AMRC https://amrc.ssec.wisc.edu/data/

## Usage
```
NAME:
   awparse - Parse AWS weather data from AMRC

USAGE:
   awparse [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     build-csv, csv  Specify data directory, glob pattern and output file name to generate a compiled CSV data file.
     help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Example
This script is very specific to the AMRC AWS weather data file
format and depends on a couple of things. Data are stored in the 
following directory structure: 
```
/dataroot/<year>/January/
                /February/
                /March/
                ...
``` 
`<year>` is a 4-digit year. And data files are expected to be named with 
the prefix: `AR######`, where `######` corresponds to 2-digit year,
2-digit month and 2-digit day, respectively. This script depends only
on the 2-digit day in the `[6:8]` position for sorting, getting
the year and month from the directory structure. Files from different
weather stations use unique naming patterns after the prefix, and a glob
pattern is used to match these. For example: 
```
awparse csv observations/ "*.100" output.csv
```
The above would consider `observations/` the data archive root
directory, would match all data files within each `<year>/<month>/` directory
using `*.100` and would generate a single sorted CSV file named
`output.csv`

### Installation
1. Clone respository
2. `cd` into repository and run `go build` to generate an `awparse` binary 