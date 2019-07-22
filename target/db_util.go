package target

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
)

func getSQLInfo(dbURL string) (string,error){
	u, err1 := url.Parse(dbURL)
	if err1!=nil {
		return "",err1
	}
	if u.Scheme != "postgres"{
		return "", fmt.Errorf("unsupported database")
	}
	host, port, err2 := net.SplitHostPort(u.Host)
	if err2!=nil {
		return "",err2
	}
	password, err3 := u.User.Password()
	if !err3  {
		return "", fmt.Errorf("no password provided")
	}

	if len(u.Path) <1 {
		return "", fmt.Errorf("no database provided or wrong")
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, u.User.Username(), password, u.Path[1:]), nil
}
//GetDatabase returns DB implementation
func GetDatabase(dbURL string) (*DatabaseImplementation,error) {
	fmt.Println("Database URL:",dbURL)
	sqlInfo,err1:=getSQLInfo(dbURL)
	if err1 != nil {
		return nil,err1
	}
	db, err4 := sql.Open("postgres", sqlInfo)
	if err4 != nil {
		return nil,err4
	}
	err5 := db.Ping()
	if err5 != nil {
		return nil,err5
	}
	impl := &DatabaseImplementation{
		mq:Postgres{},
		DB:db,
	}
	return impl,nil

}
