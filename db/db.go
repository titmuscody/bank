package db

import (
    "fmt"
    "time"
    "strconv"
    "strings"
    "crypto/hmac"
    "crypto/sha512"
    "encoding/hex"
    "math/rand"
    "database/sql"
    _ "github.com/lib/pq"
    )

var database = newPool()

func GetData() int {
    //conn := database.Get()
    //defer conn.Close()
    //val, _ := redis.Int(conn.Do("get", "temp"))
    return 12345
}


func newPool() *sql.DB {
    cStr := "user=tisourit password=3nter dbname=bank"
    db, err := sql.Open("postgres", cStr)
    if err != nil {
        panic(err)
    }
    err = db.Ping()
    if err != nil {
        panic(err)
    }
    return db
}

func makeUserReference(username string) string {
    return strings.Join([]string{"user", username}, ":")
}

func GetUserKey(username string) string {
    var key string
    database.QueryRow("select key from users where username=$1", username).Scan(&key)
    return key
}

func GetUserHash(username string) string {
    var myKey string
    var myPass string
    database.QueryRow("select key, password from users where username=$1", username).Scan(&myKey, &myPass)
    hash := hmac.New(sha512.New, []byte(myKey))
    hash.Write([]byte(myPass))
    myHash := hex.EncodeToString(hash.Sum(nil))
    stmt, err := database.Prepare("update users set key=$1 where username=$2")
    if err != nil {
        panic("key wasn't updated")
    }
    _, err = stmt.Exec(rand.Int(), username)
    if err != nil {
        panic(err)
    }
    return myHash
}

func CreateSessionId(username string) string {
    stmt, err := database.Prepare("UPDATE sessions SET id=$1, expire=$2 WHERE owner IN (SELECT id FROM users WHERE username=$3)")
    if err != nil {
        panic(err)
    }
    newId := rand.Int()
    now := time.Now().UTC()
    fmt.Println(now)
    _, err = stmt.Exec(newId, now, username)
    if err != nil {
        panic(err)
    }
    return strconv.Itoa(newId)
}

func Validate(id string) string {
    var username string
    var expire time.Time
    database.QueryRow("select users.username, sessions.expire from users, sessions where users.id=sessions.owner and sessions.id=$1;", id).Scan(&username, &expire)
    t := time.Since(expire)
    fmt.Println(time.Now().UTC())
    fmt.Println(t.Minutes())
    if t.Minutes() > 15 {
        return ""
    }
    return username
}

func refreshUserLogin(username string) {
    stmt, err := database.Prepare("UPDATE sessions SET expire=$1 WHERE owner IN (SELECT id FROM users WHERE username=$2);")
    if err != nil {
        panic(err)
    }
    _, err = stmt.Exec(time.Now().UTC(), username)
    if err != nil {
        panic(err)
    }


}
