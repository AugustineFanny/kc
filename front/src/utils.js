import router from './router'
import Vue from 'vue'

export default {
	handle: function(it, response) {
		if(response.data.code == "100103") {
			router.push("/login");
		} else if(response.data.code == "100200") {
			return response.data.data;
		} else {
			it.$message.error(response.data.msg);
		}
		return null;
	},
	dateShow: function(time) {
		if(time.slice(0, 10) == "0001-01-01")
			return ""
		else
			return time.slice(0, 10);
	},
	timeShow: function(time) {
		if(time.slice(0, 10) == "0001-01-01")
			return ""
		else
			return time.slice(0, 10) + " " + time.slice(11, 19);
	},
	amountShow: function(amount) {
		return Math.floor(amount * 100000000) / 100000000
	}
};