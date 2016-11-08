package etcd

import(
    "fmt"
    "testing"
    "net/http"
    "net/http/httptest"
    "strings"
    cli "github.com/coreos/etcd/client"
)


//==============================================================================
//                                    Login
//==============================================================================
// Login with endpoint
func Test_Login_WithEndPoint(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var db etcdDB
    db.Cfg = cli.Config{
        Endpoints: []string{ts.URL},
    }

    err := db.Login()
    if err != nil {
        T.FailNow()
    }
}


// Login without endpoint 
func Test_Login_WithoutEndPoint(T *testing.T) {
	var db etcdDB
	db.Cfg = cli.Config{
		Endpoints: []string{},
	}

	err := db.Login()
	if err == nil {
		T.FailNow()
	}
}


//============================================================================
//                                    Setup
//============================================================================
// Setup 
func Test_Setup_ValidConfig(T *testing.T){
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var db etcdDB
    err := db.Setup(ts.URL)

    if err != nil{
        T.FailNow()
    }
}

// Setup 
func Test_Setup_ServerNotFound(T *testing.T){
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(404)
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var db etcdDB
    err := db.Setup(ts.URL)

    if err == nil{
        T.FailNow()
    }
}

// Setup with no URL
func Test_Setup_NoURL(T *testing.T){
    var db etcdDB
    err := db.Setup("")

    if err == nil{
        T.FailNow()
    }
}


//============================================================================
//                                    IsSetup
//============================================================================
// IsSetup after Setup
func Test_IsSetup_After(T *testing.T){
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var db etcdDB

    err := db.Setup(ts.URL)
    if err != nil{
        T.SkipNow()
    }

    ret := db.IsSetup()
    if ret != true{
        T.FailNow()
    }
}


// IsSetup before Setup
func Test_IsSetup_Before(T *testing.T) {
    var db etcdDB

    ret := db.IsSetup()
    if ret == true {
        T.FailNow()
    }
}


//============================================================================
//                                    Set
//============================================================================
func Test_Set(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.Set("key", "value")
    if err != nil {
        T.FailNow()
    }
}


//============================================================================
//                                    Get
//============================================================================
// Get
func Test_Get_NonExistentKey(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if (r.Method == "GET"){
                w.WriteHeader(404)
            }
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    _, err = db.Get("key_non_existent")
    if err == nil {
        T.FailNow()
    }
}


//============================================================================
//                                    IsDir
//============================================================================
// IsDir
func Test_IsDir_NotDir(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if (r.Method == "GET"){
                w.WriteHeader(404)
            }
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var is_dir bool
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err, is_dir = db.IsDir("key_IsDir")
    if err != nil {
        T.SkipNow()
    }
    if is_dir == true {
        T.FailNow()
    }
}


//============================================================================
//                                    IsKey
//============================================================================
// IsKey
func Test_IsKey_ExistentKey(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var is_key bool
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    is_key, err = db.IsKey("key")
    if err != nil {
        T.SkipNow()
    }
    if is_key != true {
        T.FailNow()
    }
}


// IsKey
func Test_IsKey_NonExistentKey(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if (r.Method == "GET"){
                w.WriteHeader(404)
            }
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var is_key bool
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    is_key, err = db.IsKey("section_IsKey")
    if err != nil {
        T.SkipNow()
    }
    if is_key == true {
        T.FailNow()
    }
}


//============================================================================
//                                    Update
//============================================================================
// Update
func Test_Update(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.Update("key", "value", true)
    if err != nil {
        T.FailNow()
    }
}


//============================================================================
//                                    Del
//============================================================================
// Del (Existing key)
func Test_Del_ExistentKey(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.Del("key")
    if err != nil {
        T.FailNow()
    }
}


// Del (Non Existing key)
func Test_Del_NonExistentKey(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if (r.Method == "DELETE"){
                w.WriteHeader(404)
            }
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.Del("key")
    if err == nil {
        T.FailNow()
    }
}


//============================================================================
//                                    CreateSection
//============================================================================
// CreateSection (non-existing)
func Test_CreateSection_NonExistent(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.CreateSection("section")
    if err != nil {
        T.FailNow()
    }
}

// CreateSection (Existing)
func Test_CreateSection_Existent(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if ((r.Method == "PUT") && (strings.Contains(r.URL.Path,"section"))){
                w.WriteHeader(404)
            }
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.CreateSection("section")
    if err == nil {
        T.FailNow()
    }
}


//============================================================================
//                                    DeleteSection
//============================================================================
// DeleteSection (Existing)
func Test_DeleteSection_Existent(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.DeleteSection("section")

    if err != nil{
        T.FailNow()
    }
}

// DeleteSection (Non-Existing)
func Test_DeleteSection_NonExistent(T *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if (r.Method == "DELETE") {
                w.WriteHeader(404)
            }
            fmt.Fprintln(w, "{}")
    }))

    defer ts.Close()

    var err error
    var db etcdDB

    err = db.Setup(ts.URL)
    if err != nil {
        T.SkipNow()
    }

    err = db.DeleteSection("section")

    if err == nil {
        T.FailNow()
    }
}
