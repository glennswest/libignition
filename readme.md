# libignition - A simple ignition library
## Description
Designed to parse and manipulate json derived igntion files in windows or linux enviornments. Inspired by 
RHCoreOS igntion configuration utilities.

### Functions
#### func Parse_ignition_string(tc string) int
Takes a ignition string, and extracts the "Storage" section, and implementes the creation of files. 

Currently supports:
  1. Inline Files 
  2. Inline Base64 Files 
  3. Remote Files 

#### func Parse_ignition_file(thefilepath string) int
Process a ignition file, using supplied path. Uses Parse_igntion_string.
