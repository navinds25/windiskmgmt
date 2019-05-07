
#  How to delete duplicate files

 - [ ] Add all the files on the filesystem to the database with the file size as the key. - Requires update. operation. The value should be struct containing FullFilePath and empty Field for CheckSum.

 - [ ] For each size, if the number of entries is One continue. Else:

 - [ ] Read the files and get the checksum from DB. ->channel: add files to channel when reading.

 - [ ] If the checksum is empty. Get the checksum and update the checksum in the database.
 - [ ] Add a checksum in a map as the key with the File details as the value. <-close channel
 - [ ] Loop over the keys of the map, if the key has only one entry continue.

 - [ ] If the key has multiple entries, check if Folders are in High Priority List Folders. -> channel
 - [ ] Check if files are in Low Priority List Files.
? Add value for priority for file to be deleted.

 - [ ] If file needs to be deleted add it to a slice filesToBeDeleted. ->channel
 - [ ] Remove the list of files from the FileSize DB.
 - [ ] Delete the file from Disk
 - [ ] Add the file to Deleted DB.
