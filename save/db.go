package save

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"main/app"
	"main/tools"

	_ "github.com/go-sql-driver/mysql"
)

// Database ..
var Database *sql.DB

// SetDatabase ..
func SetDatabase() {
	db, err := sql.Open("mysql", app.Connstr)
	if err != nil {
		panic(err)
	}
	Database = db
}

// GetAddrs 根据film ID获取播放地址列表
func GetAddrs(id int) (addrs []app.Addr) {
	s := "SELECT * FROM `addr` WHERE fid = ?"
	addrs = make([]app.Addr, 0)
	rows, err := Database.Query(s, id)
	if err != nil {
		tools.Log(err.Error())
		return
	}
	for rows.Next() {
		var addr app.Addr
		if err := rows.Scan(&addr.ID, &addr.URL, &addr.Remark, &addr.FID); err != nil {
			tools.Log(err.Error())
			continue
		}
		addrs = append(addrs, addr)
	}
	return
}

// GetFilm 获取一条完整Film的数据
func GetFilm(id int) (film app.Film) {
	s := "SELECT * FROM info WHERE id = ?"
	rows, err := Database.Query(s, id)
	if err != nil {
		return
	}
	rows.Next()
	if err = rows.Scan(&film.ID, &film.Name, &film.Director, &film.Actor, &film.Type, &film.Area, &film.Lang, &film.Year, &film.Summary, &film.Img, &film.Site, &film.Amd, &film.Stat); err != nil {
		rows.Close()
		return
	}
	rows.Close()
	addrs := GetAddrs(id)
	if len(addrs) == 0 {
		tools.Log("获取播放地址时出错：电影名：" + film.Name)
	}
	film.Addrs = addrs
	return
}

// SomeInfo 根据传入电影名查询数据库电影来源网站和地址MD5，没有返回空字符串
func SomeInfo(name string) (id int, site, md string) {
	s := "SELECT id,site,amd FROM info WHERE name = ? LIMIT 0,1"
	rows, err := Database.Query(s, name)
	if err != nil {
		tools.Log(err.Error() + ": " + name)
		return
	}
	rows.Next()
	rows.Scan(&id, &site, &md)
	rows.Close()
	return
}

// SetStat 更新stat
func SetStat(id, stat int) {
	s := "UPDATE info SET stat=? WHERE id = ?"
	st, err := Database.Prepare(s)
	if err != nil {
		tools.Log(err.Error())
		return
	}
	st.Exec(stat, id)
	st.Close()
}

// SetAmd 更新Amd
func SetAmd(id int, amd string) {
	s := "UPDATE info SET amd=? WHERE id = ?"
	st, err := Database.Prepare(s)
	if err != nil {
		tools.Log(err.Error())
		return
	}
	st.Exec(amd, id)
	st.Close()
}

// InsertFilm 插入一条完整的数据
func InsertFilm(film *app.Film) bool {
	s := "INSERT INTO info VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := Database.Prepare(s)
	if err != nil {
		tools.Log(err.Error() + ": " + film.Name)
		return false
	}
	data, err := stmt.Exec(film.ID, film.Name, film.Director, film.Actor, film.Type, film.Area, film.Lang, film.Year, film.Summary, film.Img, film.Site, film.Amd, film.Stat)
	if err != nil {
		tools.Log(err.Error() + ": " + film.Name)
		stmt.Close()
		return false
	}
	fid, err := data.LastInsertId()
	if err != nil {
		tools.Log(err.Error() + ": " + film.Name)
		stmt.Close()
		return false
	}
	stmt.Close()
	s2 := "INSERT INTO `addr`(URL,remark,fid) VALUES(?,?,?)"
	st, err := Database.Prepare(s2)
	if err != nil {
		tools.Log(err.Error() + ": " + film.Name)
		return false
	}
	for _, v := range film.Addrs {
		st.Exec(v.URL, v.Remark, fid)
	}
	st.Close()
	return true
}

// UpdateAddr 更新播放地址
func UpdateAddr(film *app.Film) bool {
	s := "DELETE FROM `addr` WHERE fid = ?"
	stmt, err := Database.Prepare(s)
	if err != nil {
		tools.Log(err.Error() + ": " + film.Name)
		return false
	}
	_, err = stmt.Exec(film.ID)
	if err != nil {
		tools.Log(err.Error() + ": " + film.Name)
		stmt.Close()
		return false
	}
	stmt.Close()
	s2 := "INSERT INTO `addr`(URL,remark,fid) VALUES(?,?,?)"
	st, err := Database.Prepare(s2)
	if err != nil {
		tools.Log(err.Error() + ": " + film.Name)
		return false
	}
	defer st.Close()
	for _, v := range film.Addrs {
		_, err = st.Exec(v.URL, v.Remark, film.ID)
		if err != nil {
			return false
		}
	}
	return true
}

// HasData 查询表内是否有数据
func HasData() bool {
	s := "SELECT COUNT(*) FROM info"
	rows, err := Database.Query(s)
	if err != nil {
		tools.Log(err.Error())
		return false
	}
	count := 0
	rows.Next()
	if err := rows.Scan(&count); err != nil {
		tools.Log(err.Error())
		rows.Close()
		return false
	}
	rows.Close()
	return count != 0
}

//GetTmp 取出film表所有stat为1的数据的ID
func GetTmp() []int {
	var arr []int
	s := "SELECT id FROM info WHERE stat = 1"
	rows, err := Database.Query(s)
	if err != nil {
		tools.Log(err.Error())
		return arr
	}
	for rows.Next() {
		id := 0
		if err := rows.Scan(&id); err != nil {
			continue
		}
		if id != 0 {
			arr = append(arr, id)
		}
	}
	return arr
}

func imgMD5(bs []byte, url string) {
	s := "INSERT INTO img VALUES(?,?)"
	stmt, err := Database.Prepare(s)
	if err != nil {
		tools.Log("读取图片MD5错误：" + err.Error())
		return
	}
	defer stmt.Close()
	stmt.Exec(fmt.Sprintf("%X", md5.Sum(bs)), url)
}

func hasImg(bs []byte) string {
	s := "SELECT path FROM img WHERE md5 = ?"
	row, err := Database.Query(s, fmt.Sprintf("%X", md5.Sum(bs)))
	if err != nil {
		tools.Log("查找图片MD5错误：" + err.Error())
		return ""
	}
	p := ""
	row.Next()
	row.Scan(&p)
	row.Close()
	return p
}
