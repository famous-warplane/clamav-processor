# ClamAV Processor

A simple application to run ClamAV's clamdscan daemon to seperate quarantined files from good files....

## Running:
The application takes three command line flags:
* `-i`, `--inputPath` : The directory containing files to be scanned
* `-q`, `--quarantinePath` : The directory to which quarantined files will be moved
* `-o`, `--outputPath` : The directory to which safe files will be moved
* `-p`, `--period` : The interval in seconds at which the scan will run
### Example:

```shell script
cap -i /var/in -o /var/out -q /var/quarantine -p 30
```

## TODO:

- [x] Add logging
- [x] Add timer element so it can run periodically
- [] Test it properly with lots of files (virus test files can be found here https://www.eicar.org/?page_id=3950)

