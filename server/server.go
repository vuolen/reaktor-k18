package server

// type serverContext struct {
// 	*sqlx.DB
// }
//
// type server struct {
// 	context serverContext
// }
//
// func New(dbFile string) (server, error) {
// 	tdb, err := db.OpenDatabase(dbFile)
// 	if err != nil {
// 		return server{}, err
// 	}
// 	return server{serverContext{tdb}}, nil
// }
//
// func (s server) Run() error {
// 	r := mux.NewRouter()
// 	//r.Methods("GET").Path("/locations").Handler(Logger(s.handler.GetLocations, "GetLocations"))
// 	//r.Methods("GET").Path("/logs").Handler(Logger(s.handler.GetLogs, "GetLogs"))
// 	//r.Methods("GET").PathPrefix("/").Handler(http.FileServer("./public-html"))
// 	return http.ListenAndServe(":8080", r)
// }
//
// func (s server) Close() error {
// 	return s.handler.tdb.Close()
// }
