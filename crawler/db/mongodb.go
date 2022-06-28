package mongodb

import (
	"fmt"
	"os"
	"time"
	"math/rand"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{}) //   .JSONFormatter{}
	logrus.SetLevel(logrus.DebugLevel)
}

func getDB() *mgo.Database {
	mng_host := os.Getenv("MONGODB_HOST")
	mng_port := os.Getenv("MONGODB_PORT")
	mng_username := os.Getenv("MONGODB_USERNAME")
	mng_passwd := os.Getenv("MONGODB_PASSWD")

	session, err := mgo.Dial(mng_host + ":" + mng_port) // host:port
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	db := session.DB("admin")          //root user is created in the admin authentication database and given the role of root.
	db.Login(mng_username, mng_passwd) // login with user,password
	db = session.DB("fino_data")       // Change to working database
	return db
}

// Upsert mean insert record if data is not existed, othewise do update
func Upsert_data(data_field_list []Data_field) {
	var mydb = getDB()
	mydb.Login("admin", "6lotus")
	tms_collection := mydb.C("times_search")	//放入times_search collection
	for row_no, record := range data_field_list {
		filter := bson.D{{"pid", record.Pid}}
		update := bson.D{{"$set", record}}
		// opts := options.Update().SetUpsert(true)
		// result, err := tms_collection.UpdateOne(context.TODO(), filter, update, opts)
		_, err := tms_collection.Upsert(filter, update)
		if err != nil {
			logrus.Error("mongodb access error", row_no, record, err)
			panic(err)
		}
	}

	fmt.Println("data upsert done.")

	// err := tms_collection.Insert(&User{Name: "江阿迪", Age: 58})
	// if err != nil {
	// 	fmt.Println(err) // <nil>
	// 	panic(err)
	// }
	// result := []User{}
	// var result []bson.M
	// err = c.Find(bson.M{"name": "江阿迪"}).All(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for record, elem := range result {
	// 	// fmt.Println(record, elem["name"], elem["age"], elem["tel"])
	// 	logrus.Info("fetching data result:", record, elem["name"], elem["age"], elem["tel"])
	// }
}

func Get_times_press() (string, string, string, time.Time, string){
	var mydb = getDB()
	mydb.Login("admin", "6lotus")
	usr_collection := mydb.C("linebot_access_control")	
	var usr_result map[string]interface{}
	err := usr_collection.Find(nil).Sort("+push_date").One(&usr_result)
	if err != nil {
		logrus.Error(err)
	}

	push_sourcee_id := usr_result["user_id"].(string)

	sys_id := usr_result["_id"].(bson.ObjectId)
	filter := bson.D{{"_id", sys_id}}
	update := bson.D{{"$set", bson.D{{"push_date", time.Now()}}}}

	usr_collection_upt := mydb.C("linebot_access_control")	
	err = usr_collection_upt.Update(filter, update)				// 更新 有發送的日期

	tms_collection := mydb.C("times_search")	
	var tms_result []map[string]interface{}
	// err = tms_collection.Find(bson.M{"search_key": "*sys-即時新聞"}).Select(bson.M{"pub_date": 1}).Sort("-post_date").Limit(5).All(&tms_result)	
	err = tms_collection.Find(bson.M{"search_key": "*sys-即時新聞"}).Sort("-pub_date").Limit(5).All(&tms_result)
	// logrus.Info("tms_result:", tms_result)
	push_title := ""
	push_link := ""
	push_date := time.Now()
	push_media := ""
	rand.Seed(time.Now().UnixNano())  // 讓亂數不會重複
	choose_record := rand.Intn(5)	  // 產生 0~4 整數
	logrus.Info("choose:", choose_record)

	if tms_result[choose_record]["title"]!=nil && tms_result[choose_record]["page_link"]!=nil && tms_result[choose_record]["pub_date"]!=nil && tms_result[choose_record]["channel_url"]!=nil{
			push_title = tms_result[choose_record]["title"].(string)
			push_link = tms_result[choose_record]["page_link"].(string)
			push_media = tms_result[choose_record]["channel_url"].(string)
			push_date = tms_result[choose_record]["pub_date"].(time.Time)		
	}

	return push_sourcee_id, push_title, push_link, push_date, push_media	
}
