package models

import(
	"time"
	"strings"
	"strconv"
	"os"
	"path"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
)
//分类
type Category struct {
	Id 				int64
	Title 			string
	Created 		time.Time `orm:"index"`
	Views 			int64 	  `orm:"index"`
	TopicTime 		time.Time `orm:"index"`
	TopicCount 		int64
	TopicLastUserId int64
}
//文章
type Topic struct {
	Id 				int64
	Uid 			int64
	Title 			string
	Content 		string `orm:"index"`//表示对此字段添加索引
	Category		string
	Labels			string
	Attachment		string
	Created 		time.Time `orm:"index"`
	Updated 		time.Time `orm:"index"`
	Views 			int64 `orm:"index"`
	Author 			string
	ReplyTime 		time.Time `orm:"index"`
	ReplyCount 	int64
	ReplylastUserId int64
}
//评论
type Comment struct {
	Id int64
	Tid int64
	Name string
	Content string `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}


func AddReply(tid,nickname,content string) error {
	tidNum ,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}
	reply := &Comment {
		Tid: 	tidNum,
		Name:	nickname,
		Content:content,
		Created:time.Now(),
	}
	o := orm.NewOrm()
	_,err = o.Insert(reply)
	if err != nil {
		return err
	}
	//更新回复数和回复时间
	topic := &Topic{
		Id: tidNum,
	}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_,err = o.Update(topic)
	}


	return err

}
func DeleteReply(rid string) error {
	ridNum ,err := strconv.ParseInt(rid,10,64)
	if err != nil {
		return err
	}
	var tidNum int64
	o := orm.NewOrm()
	reply := &Comment{Id:ridNum}
	if o.Read(reply) == nil {
		tidNum = reply.Tid
		_,err = o.Delete(reply)
		if err != nil {
			return err
		}
	}
	//保存所有回复,采取精确统计
	replies := make([]*Comment,0)
	/*
	如果删除最后回复的评论，则最后回复时间需要更改
	*/
	qs := o.QueryTable("comment")
	_,err = qs.Filter("tid",tidNum).OrderBy("-created").All(&replies)
	if err != nil {
		return err
	}
	topic := &Topic{Id:tidNum}
	if o.Read(topic) == nil {
		//更新最后回复时间
		topic.ReplyTime = replies[0].Created
		topic.ReplyCount = int64(len(replies))
		_,err = o.Update(topic)
	}
	return err
}
func GetAllReplies(tid string) (replies []*Comment,err error) {
	tidNum, err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return nil,err
	}
	replies = make([]*Comment,0)
	o := orm.NewOrm()
	qs := o.QueryTable("comment")
	_,err = qs.Filter("tid",tidNum).All(&replies)
	return replies,err
}




func AddTopic(title ,content ,category ,label ,attachment string) error {
	//处理标签
	/*
	strings.Join([beego orm],"#$")--->beego#$orm
	*/
	label = "$"+strings.Join(strings.Split(label," "),"#$")+"#"
	


	o := orm.NewOrm()
	topic := &Topic {
		Title :     title,
		Content:    content,
		Category:   category,
		Labels:	    label,
		Attachment: attachment,
		Created:    time.Now(),
		Updated:    time.Now(),
		ReplyTime:  time.Now(),
	}
	_,err := o.Insert(topic)
	if err != nil {
		return err
	}
	//更新分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title",category).One(cate)
	if err == nil {
		//如果不存在，简单的忽略更新操作
		cate.TopicCount++
		_,err = o.Update(cate)
	}
	return err
}
func GetAllTopics(cate ,label string,isDesc bool) ([]*Topic,error) {
	o := orm.NewOrm()
	topics := make([]*Topic,0)
	qs := o.QueryTable("topic")
	var err error
	//降序，从最新的时间开始
	if isDesc {
		if len(cate)>0 {
			qs = qs.Filter("category",cate)
		}
		if len(label)>0 {
			qs = qs.Filter("labels__contains","$"+label+"#")
		}



		//字段名前加 - 号，表示降序
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)

	}
	return topics, err
}
func GetTopic(tid string) (*Topic,error) {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return nil,err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id",tidNum).One(topic)
	if err != nil {
		return nil,err
	}
	topic.Views++
	_,err = o.Update(topic)
	//还原labels为原先的形式，以空格分开
	topic.Labels = strings.Replace(strings.Replace(topic.Labels,"#"," ",-1),"$","",-1)
	return topic,err
	
}

func ModifyTopic(tid,title,category,content,label,attachment string) error {
	tidNum,err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}
	//处理标签
	label = "$"+strings.Join(strings.Split(label," "),"#$")+"#"
	var oldCate,oldAttach string
	o := orm.NewOrm()
	topic := &Topic{
		Id: tidNum,
	}
	if o.Read(topic) == nil {
		//保存旧的分类名称
		oldCate = topic.Category
		oldAttach = topic.Attachment
		topic.Title = title
		topic.Content = content
		topic.Labels = label
		topic.Attachment = attachment
		//赋值新的分类名称
		topic.Category = category
		topic.Updated = time.Now()
		_,err = o.Update(topic)
		if err != nil {
			return err
		}

	}

	//更新旧的分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title",oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_,err = o.Update(cate)
		}
	}
	//刪除旧的附件
	if len(oldAttach) > 0 {
		os.Remove(path.Join("attachment",oldAttach))
	}



	//更新新的分类统计
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title",category).One(cate)
	if err == nil {
		cate.TopicCount++
		_,err = o.Update(cate)
	}
	return nil
}
func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid,10,64)
	if err != nil {
		return err
	}
	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum,}
	if o.Read(topic)==nil {
		oldCate = topic.Category
		_,err = o.Delete(topic)
		if err != nil {
			return err
		}
	}
	if len(oldCate)>0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title",oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	return err
}





func RegisterDB(){
	//建表，数据库要存在，不然报错；
	//mysql_url 在配置文件中
	orm.RegisterDriver("mysql",orm.DRMySQL)
	orm.RegisterDataBase("default","mysql",beego.AppConfig.String("mysql_url"),10)
	
	orm.RegisterModel(new(Category),new(Topic),new(Comment))
}
//插入分类
func AddCategory(name string) error {
	db, err := sql.Open("mysql",beego.AppConfig.String("mysql_url"))
	if err != nil {
		beego.Error(err)
	}
	result, err := db.Exec("insert into category (title,created,topic_time)values(?,now(),now());",name)
	if result != nil {
		beego.Error(err)
	}
	return err
	
}
//删除分类
func DelCategory(id string) error {
	db, err := sql.Open("mysql",beego.AppConfig.String("mysql_url"))
	if err != nil {
		beego.Error(err)
	}
	id1 ,err := strconv.ParseInt(id,10,64)
	result, err := db.Exec("delete from category where id = ?",id1)
	if result != nil {
		beego.Error(err)
	}
	return err
}
//获取分类
func GetAllCategories() ([]*Category,error) {
	
	o := orm.NewOrm()
	cates := make([]*Category,0)
	qs := o.QueryTable("category")
	
	_,err := qs.All(&cates)
	
	return cates,err
}