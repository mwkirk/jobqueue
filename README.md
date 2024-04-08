# Job Queue Exercise

## Notes
* Requires Go 1.22+. Using new routing enhancements, so no gorillamux needed

## TODO
* Finish handler tests, clean up and improve existing ones
* Need tests for JobQueue type
* More logging, especially for errors in the handlers
* Check and set application/json content type headers
* Documentation and comments
* TLS skipped due to time. Auth also assumed not a requirement for this
* Configurable address/port
* DRY up http handlers
* Determine whether API error messages should be returned as JSON or text 
* Revisit errors. Some error types could be beneficial.
* Check/improve code organization, naming, package structure
* In general, it needs a thorough proof reading and clean up
