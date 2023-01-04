package snowflake

import (
	"time"

	"github.com/spf13/viper"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init() (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", viper.GetString("snowflake.startTIme"))
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(viper.GetInt64("snowflake.machineID"))
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}

//func main() {
//	if err := Init(); err != nil {
//		fmt.Printf("init failed, err:%v\n", err)
//		return
//	}
//	id := GenID()
//	fmt.Println(id)
//}
