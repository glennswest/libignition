package main;


import "os"
import "fmt"

func main(){
  args := os.Args[1:]
  mode := args[0];
  switch(mode){
     case "a": // Add file
          aname := args[1]
          fname := args[2]
          dpath := "/" + fname
          Add_base64_file(aname, fname, "", dpath)
          break
     }
}

