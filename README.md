# GoSimpleSQL
这是一个简单的sql生成器，可以避免手写sql字符串，2016年年底开始编写，团队内部使用，已经用于生产环境。


# Demo
sql查询的例子

	sSql := NewCommand(TABLE_UserEntity).
		Select(ALL).
		Where(
		MoUserEntity_Username, EQ, PARAM, AND,
		MoUserEntity_Password, EQ, PARAM,
	).Args("mhj", "123")
	
调用 sSql.GenerateCommand() 产生一下sql语句

    GenerateCommand: select  *  from User_Entity where Username=? and Password=?    

执行使用方法

    rows, err := db.Query(sSql.GenerateCommand(),sSql.GetArgs()...)
    

# Demo beego orm
结合beego的orm使用的例子
        
    func (*daoUserEntity)SelectByUserAndPassword(username, password string) *mhj_models.UserEntity {
        sql := NewCommand(TABLE_UserEntity).
            Select(ALL).
            Where(
            MoUserEntity_Username, EQ, PARAM, AND,
            MoUserEntity_Password, EQ, PARAM,
        ).Args(username, password)
    
        user := UserEntity{}
        if err := baseSession.Query(sql).QueryRow(&user); err == nil {
            return &user
        } else {
            log4go.Error(err)
            return nil
        }
    
    }
