package main
//go:generate go-bindata -o ./bindata.go data/...

import (
  "text/template"
  htemplate "html/template"
  "net/http"
  "strings"
  "fmt"
  "net"
  "strconv"
  "os"
)
var data_header, _ = Asset("data/header.html")
var data_footer, _ = Asset("data/footer.html")
var data_home, _ = Asset("data/home.html")
var data_profile, _ = Asset("data/profile.html")
var data_banktransfer, _ = Asset("data/banktransfer.html")

//DATABASE
type User struct {
	Ip string
	Name string
	Money float64
	Url string
}
var d = make(map[string]*User)

func strip_ip (addr string) string {
	ip := addr
	ip_arr := strings.Split(ip, ":")
	ip = strings.Join(ip_arr[:len(ip_arr)-1], ":")
	return ip
}

func get_userdata(ip string) User {
	if _, ok := d[ip]; ok {
	} else {
		d[ip] = &User{ip,"Unknown",500,""}
	}
	d[ip].Url = Url //global var to userdata
	return *d[ip]
}

func change_money(ip string, amount float64) {
	d[ip].Money+=-amount//user
}

//RESPONSES
func home (ip string, w http.ResponseWriter) {
	userdata := get_userdata(ip)
	tmpl, err := htemplate.New("home").Parse(string(data_home)+string(data_footer))
	if err != nil { panic(err) }
	err = tmpl.Execute(w, userdata)
	if err != nil { panic(err) }
}

func profile (ip string, path string, r *http.Request, w http.ResponseWriter) {
	if r.Method == "POST" {
		d[ip].Name=r.Form["name"][0]
	}
	userdata := get_userdata(ip)
	if test := strings.Split(path, "/"); test[1] == "safe" {
		tmpl, err := htemplate.New("xss").Parse(string(data_header)+string(data_profile)+string(data_footer))
		if err != nil { panic(err) }
		err = tmpl.Execute(w, userdata)
		if err != nil { panic(err) }
	} else {
		tmpl, err := template.New("xss").Parse(string(data_header)+string(data_profile)+string(data_footer))
		if err != nil { panic(err) }
		err = tmpl.Execute(w, userdata)
		if err != nil { panic(err) }
	}
}

func banktransfer (ip string, path string, r *http.Request, w http.ResponseWriter) {
	if r.Method == "POST" {
		if s, err := strconv.ParseFloat(r.Form["amount"][0], 64); err == nil {
			if s > 0 {
				if t := strings.Split(path, "/"); t[1] == "safe" {
					for _, element := range strings.Split(Url,"; ") {
						if strings.HasSuffix(r.Header.Get("Referer"), "://"+element+path) {
							change_money(ip,s)
						}
					}
					if strings.HasSuffix(r.Header.Get("Referer"), "://localhost"+Port+path) {
						change_money(ip,s)
					}

				} else {
					d[ip].Money-=s//user
				}
			}
		}
	}
	userdata := get_userdata(ip)
	tmpl, err := htemplate.New("banktransfer").Parse(string(data_header)+string(data_banktransfer)+string(data_footer))
	if err != nil { panic(err) }
	err = tmpl.Execute(w, userdata)
	if err != nil { panic(err) }
}

//GENERAL RESPONSE
func respond(w http.ResponseWriter, r *http.Request) {
  r.ParseForm() //needs to be done to make params accesible on post call

  path := r.URL.Path //get path of request
  path = strings.TrimSuffix(path,"/") //trim pre and suffix from path

  ip := strip_ip(r.RemoteAddr)

  //fmt.Printf("%+v\n", r)
  switch path {
	case "/safe/xss/profile":
		profile(ip, path, r, w)
	case "/vulnerable/xss/profile":
		profile(ip, path, r, w)
	case "/safe/csrf/banktransfer":
		banktransfer(ip, path, r, w)
	case "/vulnerable/csrf/banktransfer":
		banktransfer(ip, path, r, w)
	case "/home":
		home(ip, w)
	case "":
		home(ip, w)
	default:
		http.NotFound(w, r)
  }
}

//MAIN
var Url string =""
var Port string = ":8090"

func main() {
	http.HandleFunc("/", respond)

	if len(os.Args) > 1 {
		Port = ":"+os.Args[1]
	} else {
		fmt.Println("If you want to run on a different port, add the portnumber after the command")
	}

	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
				case *net.IPNet:
				ip = v.IP
				case *net.IPAddr:
				ip = v.IP
			}
			if strings.Contains(ip.String(), ".") { //ipv4 simplification
				if !strings.Contains(ip.String(), "127.0.0.1") { //ipv4 simplification
					Url += ip.String()+Port+"; "
				}
			}
		}
	}

	fmt.Println("Starting webserver on localhost"+Port+"; "+Url)

	if err := http.ListenAndServe(Port, nil); err != nil {
		panic(err)
	}
}
