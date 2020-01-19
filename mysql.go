package conn

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	dbs            = make(map[string]Repository)
	configs        = make(map[string]*MysqlConfig)
	messageFormate = "============%s\n"
	DefDb          = ""
)

type Repository struct {
	*gorm.DB
}

type MysqlConfig struct {
	host     string
	port     string
	user     string
	password string
	database string
	rep      Repository
}

func NewMysqlConfig(host, port, user, password, database string) *MysqlConfig {
	if conf, ok := configs[database]; ok {
		return conf
	} else {
		conf := &MysqlConfig{
			host:     host,
			port:     port,
			user:     user,
			password: password,
			database: database,
		}
		if len(configs) == 0 {
			DefDb = database
		}
		configs[database] = conf
		return conf
	}
}

func (m *MysqlConfig) Conn() error {
	var (
		err error
	)
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=%s&parseTime=true",
		m.user, m.password, m.host, m.port, m.database, "Local")
	fmt.Printf(messageFormate, "mysql 开始连接")
	m.rep.DB, err = gorm.Open("mysql", dbUrl)
	m.rep.DB.SingularTable(true)
	if err != nil {
		fmt.Printf(messageFormate, "mysql 连接失败")
		return err
	}
	fmt.Printf(messageFormate, "mysql 连接成功")
	//将连接放入缓存
	dbs[m.database] = m.rep
	return nil
}

func (m *MysqlConfig) Close() error {
	return m.rep.Close()
}

func (repo *Repository) GDB(databases ...string) *Repository {
	return getRepository(repo, databases...)
}

func GDB(databases ...string) *Repository {
	var repo *Repository
	return getRepository(repo, databases...)
}

func getRepository(repo *Repository, databases ...string) *Repository {
	var database string
	var repD Repository
	if len(databases) == 0 {
		database = DefDb
	} else {
		database = databases[0]
	}
	if rep, ok := dbs[database]; !ok || rep.DB == nil {
		c := configs[database]
		conf := NewMysqlConfig(c.host, c.port, c.user, c.password, c.database)
		err := conf.Conn()
		if err != nil {
			fmt.Printf(messageFormate, "mysql 连接失败")
		}
		dbs[database] = conf.rep
		repD = conf.rep
	} else {
		repD = dbs[database]
	}
	if repo == nil {
		repo = &repD
	} else {
		if repo.DB == nil {
			repo.DB = repD.DB
		}
	}
	return repo
}

func (repo *Repository) Begin() *Repository {
	if repo == nil {
		repo.GDB()
	}
	repo.DB = repo.DB.Begin()
	return repo
}

func (repo *Repository) ConTran(rep *Repository) *Repository {
	if repo == nil {
		repo = &Repository{}
	}
	repo.DB = rep.DB
	return rep
}

/*
事务特别说明:
type ARepository struct {
	Repository //如果像这样为非指针继承, 夸方法的事务如下
}
//且夸方法的这个结构体也要是非指针继承
type BRepository struct {
	Repository
}
var (
	a ARepository
	b BRepository
)
a.GDB().Begin()
a.Insert()
b.ConTran(&a.Repository).Insert() //将事务继续往下丢执行sql
b.Commit()

//如果是指针
type CRepository struct {
	*Repository
}

type DRepository struct {
	Repository
}
var (
	c CRepository
	d DRepository
)

c.Repository = c.GDB()
c.Begin()
c.Insert()
d.ConTran(c.Repository).Insert()
d.Commit()
*/
