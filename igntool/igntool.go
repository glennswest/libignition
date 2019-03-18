package main;


import "os"
import "log"
import . "github.com/glennswest/libignition/ignition"


func main(){
  if (len(os.Args) == 1){
     help()
     return
     }
  args := os.Args[1:]
  mode := args[0];
  switch(mode){
     case "a": // Add file
          aname := args[1]
          fname := args[2]
          dpath := "/" + fname
          log.Printf("Adding %s to %s\n",fname,aname)
          Add_base64_file(aname, fname, "", dpath)
     default:
          log.Printf("igntool: Invalid command\n")
          help();

     }
}

func help(){
   log.Printf("Command Line Help\n")
   log.Printf("   igntool a ignfile filetoadd\n")
   return
}
