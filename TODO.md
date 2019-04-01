# To-do

1. Input validation and sanitisation. Currently the API assumes the payload is well-formatted.
1. The handlers should never return raw errors to the user. We should replace all error returns with error handling which logs the issue then returns a generic error (4XX/5XX to the user).
1. Authentication should be added - an API key or such.
1. The insert hander should not overwrite existing objects. Our upsert function should allow guarding against overwrite.
