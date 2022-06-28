package mongodb

import (
	"time"
)

// times_search 文件的 schema 
type Data_field struct {
	Pid		  			string    	`json:"pid"`
	Pub_date  			time.Time 	`json:"post_date"`
	Title   			string    	`json:"title"`
	Page_link 			string    	`json:"url"`
	Search_key 			string 		`json:"search_key"`
	Author 	  			string 		`json:"author"`
	Content   			string    	`json:"content"`
	Score				int			`json:"dcore"`
	Cat_tag   			string    	`json:"cat_tag"`
	Channel_url 		string		`json:"channel_url"`
	Update_datetime 	time.Time 	`json:"update_datetime"`
}


