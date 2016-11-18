package controller

import (
	. "app/common"
	logger "app/logger"
	. "app/models"
	"net/http"

	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	cache "github.com/patrickmn/go-cache"
)

const (
	PageSize int = 10
)

func AppointmentActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if aid == "" {
		r.JSON(200, Resp{1105, "添加活动失败,aid不能为空", nil})
	}
	var record Record
	record.Aid = aid
	record.Uid = uid
	record.Type = 0
	err := dbmap.Insert(&record)
	CheckErr(err, "AppointmentActivity insert failed")
	if err != nil {
		r.JSON(200, Resp{1105, "添加活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "预约成功", nil})
	}
}

func PlayActivity(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	var record Record
	record.Aid = req.PostFormValue("aid")
	record.Uid = uid
	record.Type = 1
	err := dbmap.Insert(&record)
	CheckErr(err, "PayActivity insert failed")
	r.JSON(200, Resp{0, "ok", nil})
}

func GetHomeList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	var recomendActivities []QActivity
	var activities []QActivity
	_, err := dbmap.Select(&recomendActivities, `SELECT *, (SELECT count(*) FROM t_record WHERE type = 2 AND uid = ? AND  aid= t.aid)  AS pay_state, 
	(SELECT count(*) FROM t_record WHERE type = 0 AND uid = ? AND  aid= t.aid)  AS appoint_state  FROM t_activity t
	WHERE is_recommend = 1 ORDER BY activity_state ASC, create_time DESC`, uid, uid)
	CheckErr(err, "get recomend list")
	_, err = dbmap.Select(&activities, `SELECT *, (SELECT count(*) FROM t_record WHERE type = 2 AND uid = ? AND  aid= t.aid)  AS pay_state, 
	(SELECT count(*) FROM t_record WHERE type = 0 AND uid = ? AND  aid= t.aid)  AS appoint_state  FROM t_activity t
	WHERE is_recommend = 0 ORDER BY activity_state ASC, create_time DESC LIMIT ?`, uid, uid, PageSize+1)
	CheckErr(err, "get Activity List")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		var hasmore bool
		logger.Info(len(activities))
		if len(activities) > PageSize {
			hasmore = true
			activities = activities[:PageSize]
		} else {
			hasmore = false
		}
		r.JSON(200, Resp{0, "查询活动成功", map[string]interface{}{
			"hasmore": hasmore, "recommend": recomendActivities, "general": activities}})
	}
}

func GetMoreActivityList(req *http.Request, params martini.Params, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	lastAid := params["lastAid"]
	var activity QActivity
	err := dbmap.SelectOne(&activity, "SELECT * FROM t_activity WHERE aid = ? ", lastAid)
	logger.Info("GetMoreActivityList..", activity.Created)
	var activities []QActivity
	_, err = dbmap.Select(&activities, `SELECT *, (SELECT count(*) FROM t_record WHERE type = 2 AND uid = ? AND  aid= t.aid)  AS pay_state,
	(SELECT count(*) FROM t_record WHERE type = 0 AND uid = ? AND  aid= t.aid)  AS appoint_state  
	FROM t_activity t 
	WHERE t.create_time < ? AND t.activity_state >= ? AND t.is_recommend = 0 ORDER BY t.activity_state ASC, t.create_time DESC  LIMIT ?`, uid, uid, activity.Created, activity.ActivityState, PageSize+1)
	CheckErr(err, "GetActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		var hasmore bool
		if len(activities) > PageSize {
			hasmore = true
			activities = activities[:PageSize]
		} else {
			hasmore = false
		}
		r.JSON(200, Resp{0, "查询活动成功", map[string]interface{}{"hasmore": hasmore, "general": activities}})
	}
}

func GetLiveActivityList(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	var activities []QActivity
	_, err := dbmap.Select(&activities, `SELECT *, (SELECT count(*) FROM t_record WHERE type = 2 AND uid = ? AND  aid= t.aid)  AS pay_state, 
	(SELECT count(*) FROM t_record WHERE type = 0 AND uid = ? AND  aid= t.aid)  AS appoint_state  
	FROM t_activity t WHERE t.activity_state = 1 AND t.stream_type = 0 ORDER BY t.create_time DESC`, uid, uid)
	CheckErr(err, "GetLiveActivityList select failed")
	if err != nil {
		r.JSON(200, Resp{1104, "查询活动失败", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activities})
	}
}

func GetActivityDetail(args martini.Params, r render.Render, dbmap *gorp.DbMap) {
	var activity QActivity
	err := dbmap.SelectOne(&activity, "select * from t_activity where aid =?", args["id"])
	CheckErr(err, "GetActivity select failed")
	if err != nil {
		r.JSON(200, Resp{1103, "活动不存在", nil})
	} else {
		r.JSON(200, Resp{0, "查询活动成功", activity})
	}
}

func JoinGroup(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if uid == "" || aid == "" {
		logger.Info("JoinGroup", "uid or aid is ''")
	} else {
		_, err := dbmap.Exec(`INSERT INTO t_activity_user_online VALUE(NULL,?,?,now())`, aid, uid)
		logger.Info("JoinGroup ", err)
	}
	r.JSON(200, Resp{0, "成功", nil})
}

func LeaveGroup(req *http.Request, r render.Render, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if uid == "" || aid == "" {
		logger.Info("LeaveGroup", "uid or aid is ''")
	} else {
		_, err := dbmap.Exec(`DELETE FROM t_activity_user_online WHERE aid = ? AND uid = ?`, aid, uid)
		logger.Info("LeaveGroup ", err)
	}
	r.JSON(200, Resp{0, "成功", nil})
}

func GetLiveActivityMemberCount(req *http.Request, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	req.ParseForm()
	aid := req.PostFormValue("aid")
	if uid == "" || aid == "" {
		logger.Info("LeaveGroup", "uid or aid is ''")
		r.JSON(200, Resp{0, "缺少参数", nil})
	} else {
		count, err := dbmap.SelectInt("SELECT COUNT(*) FROM t_activity_user_online WHERE aid = ?", aid)
		if err == nil {
			r.JSON(200, Resp{0, "获取在线成员信息成功", map[string]int{"count": int(count)}})
		} else {
			r.JSON(200, Resp{1402, "获取在线成员信息失败", nil})
		}
	}
}

func GetLiveActivityMemberList(req *http.Request, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	uid := req.Header.Get("uid")
	query := req.URL.Query()
	var aid string
	if len(query["aid"]) > 0 {
		aid = query["aid"][0]
	}
	if uid == "" || aid == "" {
		logger.Info("LeaveGroup", "uid or aid is ''")
		r.JSON(200, Resp{0, "缺少参数", nil})
	} else {
		var users []OnlineUser
		_, err := dbmap.Select(&users, `SELECT o.uid,u.avatar,u.nickname FROM t_activity_user_online o LEFT JOIN t_user u ON o.uid = u.uid WHERE o.aid = ?`, aid)
		if err == nil {
			r.JSON(200, Resp{0, "获取在线成员信息成功", users})
		} else {
			r.JSON(200, Resp{1402, "获取在线成员信息失败", nil})
		}
	}
}

func GetSharePage(params martini.Params, r render.Render, c *cache.Cache, dbmap *gorp.DbMap) {
	platform := params["platform"]
	logger.Info("platform", platform)
	var activity QActivity
	err := dbmap.SelectOne(&activity, "select * from t_activity where aid =?", params["id"])
	CheckErr(err, "GetActivity select failed")
	///apple 下载地址 https://itunes.apple.com/cn/app/qq/id444934666
	//应用宝下载地址
	if err == nil {
		r.JSON(200, Resp{0, "获取成功", activity})
	} else {
		r.JSON(200, Resp{1103, "获取在线成员信息失败", nil})
	}
}
