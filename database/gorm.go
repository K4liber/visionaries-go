package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

/*Mem ...
 */
type Mem struct {
	ID             int
	Signature      string `gorm:"type:varchar(100)"`
	ImgExt         string `gorm:"column:imgExt"`
	DateTime       string `gorm:"column:dateTime"`
	AuthorNickname string `gorm:"column:authorNickname"`
	Category       string `gorm:"type:varchar(50)"`
	Points         int
	Views          int
	Like           bool
}

/*MemDB ...
 */
type MemDB struct {
	ID             int
	Signature      string `gorm:"type:varchar(100)"`
	ImgExt         string `gorm:"column:imgExt"`
	DateTime       string `gorm:"column:dateTime"`
	AuthorNickname string `gorm:"column:authorNickname"`
	Category       string `gorm:"type:varchar(50)"`
	Points         int
	Views          int
}

/*MemPoint ...
 */
type MemPoint struct {
	ID             int
	MemID          int    `gorm:"column:memId"`
	AuthorNickname string `gorm:"column:authorNickname"`
	DateTime       string `gorm:"column:dateTime"`
}

/*Comment ...
 */
type Comment struct {
	ID             int
	MemID          int    `gorm:"column:memId"`
	AuthorNickname string `gorm:"column:authorNickname"`
	AuthorPhoto    string `gorm:"column:authorPhoto"`
	Content        string
	DateTime       string `gorm:"column:dateTime"`
	Points         int
	Like           bool
}

/*CommentDB ...
 */
type CommentDB struct {
	ID             int
	MemID          int    `gorm:"column:memId"`
	AuthorNickname string `gorm:"column:authorNickname"`
	AuthorPhoto    string `gorm:"column:authorPhoto"`
	Content        string
	DateTime       string `gorm:"column:dateTime"`
	Points         int
}

/*CommentPoint ...
 */
type CommentPoint struct {
	ID             int
	CommentID      int    `gorm:"column:commentId"`
	AuthorNickname string `gorm:"column:authorNickname"`
	DateTime       string `gorm:"column:dateTime"`
	MemID          int    `gorm:"column:memId"`
}

//SetDB - set database connection
func SetDB(gormDB *gorm.DB) {
	db = gormDB
}

func init() {
	fmt.Println("Init package database")
}

/*GetMem - ...
() */
func GetMem(memID int, nickname string) Mem {
	var mem Mem
	db.Where("id = ?", memID).First(&mem)
	var views = mem.Views + 1
	var memDB MemDB
	db.Table("mem").Model(&memDB).Where("id = ?", memID).Update("views", views)
	mem.Like = GetMemLike(memID, nickname)
	return mem
}

/*GetComments - ...
() */
func GetComments(memID int, nickname string) []Comment {
	var comments []Comment
	db.Where("memId = ?", memID).Find(&comments)
	for index, comment := range comments {
		comments[index].Like = GetCommentLike(comment.ID, nickname)
	}
	return comments
}

/*GetMemLike - return first memLike from db
(memID int, authorNickname string) */
func GetMemLike(memID int, authorNickname string) bool {
	var iLikeIt = false
	var memPoint MemPoint
	db.Table("memPoint").Where("memId = ? AND authorNickname = ?", memID, authorNickname).Find(&memPoint)
	if memPoint.AuthorNickname == authorNickname {
		iLikeIt = true
	}
	return iLikeIt
}

/*GetCommentLike - return first commentLike from db
(commentID int, nickname string) */
func GetCommentLike(commentID int, nickname string) bool {
	var iLikeIt = false
	var commentPoint CommentPoint
	db.Table("commentPoint").Where("commentId = ? AND authorNickname = ?", commentID, nickname).Find(&commentPoint)
	if commentPoint.AuthorNickname == nickname {
		iLikeIt = true
	}
	return iLikeIt
}

/*GetMems - return all mems from database
(authorNickame string) - argument to verified Like column */
func GetMems(authorNickname string) []Mem {
	var mems []Mem
	db.Find(&mems)
	for index, mem := range mems {
		mems[index].Like = GetMemLike(mem.ID, authorNickname)
	}
	return mems
}

/*GetProfileMems - return user mems from database by nickname
(nickname string) */
func GetProfileMems(nickname string) []Mem {
	var mems []Mem
	db.Where("authorNickname = ?", nickname).Find(&mems)
	for index, mem := range mems {
		mems[index].Like = GetMemLike(mem.ID, nickname)
	}
	return mems
}

/*GetProfileComments - return user comments from database by nickname
(nickname string) */
func GetProfileComments(nickname string) []Comment {
	var comments []Comment
	db.Where("authorNickname = ?", nickname).Find(&comments)
	for index, comment := range comments {
		comments[index].Like = GetCommentLike(comment.ID, nickname)
	}
	return comments
}

/*GetProfileCommentLike - return user commentsLikes from database by nickname
(nickname string) */
func GetProfileCommentLike(nickname string) []CommentPoint {
	var commentPoints []CommentPoint
	db.Table("commentPoint").Where("authorNickname = ?", nickname).Find(&commentPoints)
	return commentPoints
}

/*GetProfileMemLike - return user commentsLikes from database by nickname
(nickname string) */
func GetProfileMemLike(nickname string) []MemPoint {
	var memPoints []MemPoint
	db.Table("memPoint").Where("authorNickname = ?", nickname).Find(&memPoints)
	return memPoints
}

/*GetCategoryMems - return category mems
(category string, nickname string) */
func GetCategoryMems(category string, nickname string) []Mem {
	var mems []Mem
	db.Where("category = ?", category).Find(&mems)
	for index, mem := range mems {
		mems[index].Like = GetMemLike(mem.ID, nickname)
	}
	return mems
}

/*InsertMem - insert article into DB and return his ID
(mem Mem) */
func InsertMem(mem MemDB) int {
	db.Table("mem").Create(&mem)
	return mem.ID
}

/*InsertComment - insert comment into DB and return his ID
(comment Comment) */
func InsertComment(comment CommentDB) CommentDB {
	db.Table("comment").Create(&comment)
	return comment
}

/*UserCommentUpdate - updates comments photos
(nickname, photoURL) */
func UserCommentUpdate(nickname string, photoURL string) {
	var comment CommentDB
	db.Table("comment").Model(&comment).Where("authorNickname = ?", nickname).Update("authorPhoto", photoURL)
}

/*AddMemPoint - ...
(memPoint) */
func AddMemPoint(memPoint MemPoint) MemPoint {
	var memPointDB MemPoint
	db.Table("memPoint").Where("authorNickname = ? AND memId = ?", memPoint.AuthorNickname, memPoint.MemID).Find(&memPointDB)
	if memPointDB.ID == 0 {
		db.Table("memPoint").Create(&memPoint)
	}
	return memPoint
}

/*AddCommentPoint - ...
(commentPoint) */
func AddCommentPoint(commentPoint CommentPoint) CommentPoint {
	var commentPointDB CommentPoint
	db.Table("commentPoint").Where("authorNickname = ? AND commentId = ?", commentPoint.AuthorNickname, commentPoint.CommentID).Find(&commentPointDB)
	if commentPointDB.ID == 0 {
		db.Table("commentPoint").Create(&commentPoint)
	}
	return commentPoint
}

/*DeleteMemPoint - ...
(memPoint) */
func DeleteMemPoint(memPoint MemPoint) {
	var memPointDB MemPoint
	db.Table("memPoint").Where("authorNickname = ? AND memId = ?", memPoint.AuthorNickname, memPoint.MemID).Find(&memPointDB)
	db.Table("memPoint").Delete(&memPointDB)
}

/*DeleteCommentPoint - ...
(commentPoint) */
func DeleteCommentPoint(commentPoint CommentPoint) {
	var commentPointDB CommentPoint
	db.Table("commentPoint").Where("authorNickname = ? AND commentId = ?", commentPoint.AuthorNickname, commentPoint.CommentID).Find(&commentPointDB)
	db.Table("commentPoint").Delete(&commentPointDB)
}

/*UpdateMemPoints - ...
(memID) */
func UpdateMemPoints(memID int) {
	var memPoints []MemPoint
	db.Table("memPoint").Where("memId = ?", memID).Find(&memPoints)
	var points = len(memPoints)
	var mem MemDB
	db.Table("mem").Model(&mem).Where("id = ?", memID).Update("points", points)
}

/*UpdateCommentPoints - ...
(commentID) */
func UpdateCommentPoints(commentID int) CommentDB {
	var commentPoints []CommentPoint
	db.Table("commentPoint").Where("commentId = ?", commentID).Find(&commentPoints)
	var points = len(commentPoints)
	var comment CommentDB
	db.Table("comment").Model(&comment).Where("id = ?", commentID).Update("points", points)
	db.Table("comment").Where("id = ?", commentID).First(&comment)
	return comment
}

/*DeleteMemPoints - ...
(memID) */
func DeleteMemPoints(memID int) {
	db.Table("memPoint").Where("memId = ?", memID).Delete(MemPoint{})
}

/*DeleteComments - ...
(memID) */
func DeleteComments(memID int) {
	db.Table("comment").Where("memId = ?", memID).Delete(CommentDB{})
}

/*DeleteCommentPoints - ...
(memID) */
func DeleteCommentPoints(memID int) {
	db.Table("commentPoint").Where("memId = ?", memID).Delete(CommentPoint{})
}

/*DeleteMem - ...
(memID) */
func DeleteMem(memID int) {
	db.Table("mem").Where("id = ?", memID).Delete(MemDB{})
}

/*DeleteComment - ...
(commentID) */
func DeleteComment(commentID int) {
	db.Table("comment").Where("id = ?", commentID).Delete(CommentDB{})
}

/*DeleteCommentPointsByCommID - ...
(commentID) */
func DeleteCommentPointsByCommID(commentID int) {
	db.Table("commentPoint").Where("commentId = ?", commentID).Delete(CommentPoint{})
}
