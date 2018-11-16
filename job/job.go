package job

import (
	"github.com/astaxie/beego"
)

func Start() {
	if beego.BConfig.RunMode == "dev" {
		//go UpdateBTCNowAmount()
		//go UpdateETHNowAmount()
		//go BtcdoJob()
		go MiningJob()
		//go GatheringJob()
		//go BTCJob()
		//go ETHJob()
		//go USDTJob()
		//go TokenJob()
	} else {
		//go UpdateBTCNowAmount()
		//go UpdateETHNowAmount()
		//go BtcdoJob()
		go MiningJob()
		go GatheringJob()
		//go BTCJob()
		//go ETHJob()
		//go USDTJob()
		//go TokenJob()
	}
}
