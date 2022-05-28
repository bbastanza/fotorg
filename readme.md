# Fotorg

> A file organization program written in GO

### what is it?

* fotorg takes all the files in a source directory and copys 
them to a destination directory organized in directories for each file
extension found in the source directory

### configuration

* A configuration file should be located at ~/.config/fotorg/config.json

* Properties
   * source: path to the source location of the files
   * destination: path to the destination location
   
```
{
    "source": "/home/username/Pictures/source-example",
    "destination": "/home/username/Pictures/destination-example"
}
```

### install

1. install go
2. place fotorg executable file in somewhere in your PATH
3. run fotorg
