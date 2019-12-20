package conn

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbs            = make(map[string]*Repository)
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
		if len(configs) == 0{
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
	fmt.Println(m.host)
	m.rep.DB, err = gorm.Open("mysql", dbUrl)
	m.rep.DB.SingularTable(true)
	if err != nil {
		fmt.Printf(messageFormate, "mysql 连接失败")
		return err
	}
	fmt.Printf(messageFormate, "mysql 连接成功")
	//将连接放入缓存
	dbs[m.database] = &m.rep
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
	if len(databases) == 0 {
		database = DefDb
	} else {
		database = databases[0]
	}

	if repo.DB == nil {
		if dbs[database].DB == nil {
			c := configs[database]
			conf := NewMysqlConfig(c.host, c.port, c.user, c.password, c.database)
			err := conf.Conn()
			if err != nil {
				fmt.Printf(messageFormate, "mysql 连接失败")
			}
			dbs[database] = &conf.rep
			repo = &conf.rep
		} else {
			repo = dbs[database]
		}
	}
	return repo
}