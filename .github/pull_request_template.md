## Purpose
_Describe the problem or feature in addition to a link to the issues._

## Approach
_How does this Pull Request address the problem?_

## Testing

 - There should be some unit tests for every behaviour change or new feature
 - If you're adding a new configuration option, the default configuration file
   (.gitmux.yml) should be modified with the new option
 - There's a test that verifies that gitmux output always remains the same
   across versions when using the default configuration. You can run it with:
  `go test -run TestScriptst`
 - `README.md` should always be updated.