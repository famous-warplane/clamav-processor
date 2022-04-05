# ClamAV Processor

A simple application to run ClamAV's clamdscan daemon to seperate quarantined files from good files....

## Running:
The application takes three command line arguments:
* directoryToScan
* quarantineDirectory
* safeDirectory

### Example:

```shell script
go run main.go "/dir/to/scan/*" "/quarantine/dir/" "/safe/dir/"
```

## TODO:
* Add logging
* Add timer element so it can run periodically
* Test it properly with lots of files (virus test files can be found here https://www.eicar.org/?page_id=3950)

