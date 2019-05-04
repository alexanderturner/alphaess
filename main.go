package main

// DBURL, DBUSER, DBPASS, DBNAME, ESSUSER, ESSPASS, ESSSN

func main() {
	//test ENV Vars
	// if os.Getenv("DBURL") == "" {
	// 	log.Fatal("DBURL not set")
	// }
	// if os.Getenv("DBUSER") == "" {
	// 	log.Fatal("DBUSER not set")
	// }
	// if os.Getenv("DBPASS") == "" {
	// 	log.Fatal("DBPASS not set")
	// }
	// if os.Getenv("DBNAME") == "" {
	// 	log.Fatal("DBNAME not set")
	// }
	// if os.Getenv("ESSUSER") == "" {
	// 	log.Fatal("ESSUSER not set")
	// }
	// if os.Getenv("ESSPASS") == "" {
	// 	log.Fatal("ESSPASS not set")
	// }
	// if os.Getenv("ESSSN") == "" {
	// 	log.Fatal("ESSSN not set")
	// }

	done := make(chan bool)

	a := new(App)
	a.DB = InfluxDBClient()
	a.ApiPoller()

	<-done

	// router := mux.NewRouter()

	// router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	// router.HandleFunc("/api/login", controllers.Authenticate).Methods("POST")
	// router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	// router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	// router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler
	//
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8000" //localhost
	// }
	//
	// // fmt.Println(port)
	//
	// err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	// if err != nil {
	// 	fmt.Print(err)
	// }

}
