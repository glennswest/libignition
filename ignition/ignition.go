package ignition;
  
import(
            "os"
            "fmt"
            "net/url"
            "strings"
            "io"
            "io/ioutil"
            "path/filepath"
            "encoding/base64"
            "github.com/tidwall/gjson"
            "github.com/tidwall/sjson"
            "net/http"
            "strconv"
            "log"
        )

// Examples
//{
//  "ignition": { "version": "2.2.0" },
//  "storage": {
//    "files": [{
//      "filesystem": "root",
//      "path": "/foo/bar",
//      "mode": 420,
//      "contents": { "source": "data:,example%20file%0A" }
//    }]
//  }
//}
//{
//  "ignition": { "version": "2.2.0" },
//  "storage": {
//    "files": [{
//      "filesystem": "root",
//      "path": "/foo/bar",
//      "mode": 420,
//      "contents": {
//        "source": "http://example.com/asset",
//        "verification": { "hash": "sha512-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" }
//      }
//    }]
//  }
//}
//    {
//        "filesystem" : "root",
//        "path" : "./foo/motd",
//        "mode" : 420,
//        "contents" : {
//          "source" : "data:text/plain;charset=utf-8;base64,VGhpcyBpcyB0aGUgYm9vdHN0cmFwIG5vZGU7IGl0IHdpbGwgYmUgZGVzdHJveWVkIHdoZW4gdGhlIG1hc3RlciBpcyBmdWxseSB1cC4KClRoZSBwcmltYXJ5IHNlcnZpY2UgaXMgImJvb3RrdWJlLnNlcnZpY2UiLiBUbyB3YXRjaCBpdHMgc3RhdHVzLCBydW4gZS5nLgoKICBqb3VybmFsY3RsIC1iIC1mIC11IGJvb3RrdWJlLnNlcnZpY2UK",
//          "verification" : {}
//        }
//      }



//Example:
//func main(){
//  parse_ignition_file("./ignition-test.json");
//}

const ignition_base_json string = `{
  "ignition": { "version": "2.2.0" },
  "storage": {
    "files": []
  }
}`

func file_is_exists(f string) bool {
    _, err := os.Stat(f)
    if os.IsNotExist(err) {
        return false
    }
    return err == nil
}

func New_ignition_file(path string) int {
	tdir := filepath.Dir(path)
        os.MkdirAll(tdir,os.ModePerm)
        td := []byte(ignition_base_json)
        err := ioutil.WriteFile(path, td, 0644)
        if (err != nil){
           log.Printf("Error: New_ignition_file %s Failed - %s\n",path,err)
           return(-1)
           }
        return(1)
}

func Find_storage_idx(tc string,destpath string) int {

        result := gjson.Get(tc,"storage.files");
        files := result.Array();
        for idx,tfile := range files {
            tpath := gjson.Get(tfile.String(),"path").String()
            if (tpath == destpath){
               return(idx)
               }
            }
        return(len(files)) // For append
}

func Add_base64_file(jsonpath string, filetoadd string, destfs string, destpath string) int {

        if (file_is_exists(jsonpath) == false){
           New_ignition_file(jsonpath)
           }
	jsb, err := ioutil.ReadFile(jsonpath)
	if (err != nil) {
                log.Printf("Add_base64_file: Failed to read json(%s) - %s->%s\n",err,jsonpath, filetoadd, destpath)
		return(-1)
                }
        js := string(jsb)

	content, err := ioutil.ReadFile(filetoadd)
	if (err != nil) {
                log.Printf("Add_base64_file: Failed(%s) - %s->%s\n",err,jsonpath, filetoadd, destpath)
		return(-1)
                }
       bcontent := base64.StdEncoding.EncodeToString(content)
       idx := Find_storage_idx(js,destpath)
       vname := "storage.files." + strconv.Itoa(idx)
       js,_ = sjson.Set(js,vname + ".contents.source", "data:text/plain;charset=utf-8;base64," + bcontent)
       js,_ = sjson.Set(js,vname + ".mode", 420)
       js,_ = sjson.Set(js,vname + ".filesystem", "")
       js,_ = sjson.Set(js,vname + ".path", destpath )
       d := []byte(js)
       err = ioutil.WriteFile(jsonpath, d, 0644)
       if (err != nil){
           log.Printf("Error: Add_base64_file %s Failed - %s\n",jsonpath,err)
           return(-1)
           }
       return(0)
}


func Parse_ignition_string(tc string) int {
	version := gjson.Get(tc, "ignition.version");
        if (version.String() == ""){
           fmt.Printf("Invalid file");
           return(-1);
           }
        result := gjson.Get(tc,"storage.files");
        files := result.Array();
        for _,tfile := range files {
            tpath := gjson.Get(tfile.String(),"path").String();
            tmode := gjson.Get(tfile.String(),"mode").Int();
            tdata := gjson.Get(tfile.String(),"contents.source").String();
            idx := strings.Index(tdata,":")+1;
            thetype := tdata[:idx];
            fmt.Printf("path: %s type: %s mode %o\n",tpath,thetype,tmode);
            tdir := filepath.Dir(tpath);
            fmt.Printf("%s\n",tdir);
            os.MkdirAll(tdir, os.ModePerm);
            fmt.Printf("Type: path: %s type: %s mode %o\n",tpath,thetype,tmode);
            switch thetype {
               case "data:":
                    cidx := strings.Index(tdata,",");
                    dtype := tdata[idx:cidx];
                    if (strings.Contains(dtype,"base64")){
                       dtype = "base64";
                       }
                    fmt.Printf("Dtype=%s\n",dtype);
                    switch dtype {
                        case "":
                          untc,_ := url.QueryUnescape(tdata[cidx+1:]);
                          td := []byte(untc);
                          err := ioutil.WriteFile(tpath, td, os.FileMode(tmode));
                          if (err != nil){
                             fmt.Printf("Failed to Write %s: %s\n",tpath,err);
                             }
                        case "base64":
                          untc,_ := base64.StdEncoding.DecodeString(tdata[cidx+1:]);
                          td := []byte(untc);
                          err := ioutil.WriteFile(tpath, td, os.FileMode(tmode));
                          if (err != nil){
                             fmt.Printf("Failed to Write %s: %s\n",tpath,err);
                             }
                          }
               case "http:","https:":
                    err := downloadfile(tpath,tdata);
                    if (err != nil){
                       fmt.Printf("Download Failed: %s - %s\n",tpath,err);
                       }
               default:
                  fmt.Printf("Invalid Type: path: %s type: %s mode %o\n",tpath,thetype,tmode);
               }
                 
            }
        return(0);
}

func Parse_ignition_file(thefilepath string) int {

    b, err :=ioutil.ReadFile(thefilepath);
    if err != nil {
      fmt.Print(err);
      return 0;
      }
    content := string(b);
    result := Parse_ignition_string(content);
    return(result);

}


// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadfile(filepath string, url string) error {

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}

