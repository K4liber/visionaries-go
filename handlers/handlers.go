package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	"visionaries/database"

	"github.com/gorilla/mux"
)

//HostName ...
const HostName = "localhost:3000"

//APIToken - aut0 api token
const APIToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IlFqZENNMFV5T0VNM00wRXdNVFZHUVVRMk56RkdOVGMyTkRZMk0wSXdRME00TkVVelFVUkVPUSJ9.eyJpc3MiOiJodHRwczovL2s0bGliZXIuZXUuYXV0aDAuY29tLyIsInN1YiI6ImNyeTUwV0ZYNVJrSUlVbGY1aUdiRnFwQURqQ1g3UlRrQGNsaWVudHMiLCJhdWQiOiJodHRwczovL2s0bGliZXIuZXUuYXV0aDAuY29tL2FwaS92Mi8iLCJleHAiOjE1MDAxNTIxMjYsImlhdCI6MTQ5MTUxMjEyNiwic2NvcGUiOiJyZWFkOmNsaWVudF9ncmFudHMgY3JlYXRlOmNsaWVudF9ncmFudHMgZGVsZXRlOmNsaWVudF9ncmFudHMgdXBkYXRlOmNsaWVudF9ncmFudHMgcmVhZDp1c2VycyB1cGRhdGU6dXNlcnMgZGVsZXRlOnVzZXJzIGNyZWF0ZTp1c2VycyByZWFkOnVzZXJzX2FwcF9tZXRhZGF0YSB1cGRhdGU6dXNlcnNfYXBwX21ldGFkYXRhIGRlbGV0ZTp1c2Vyc19hcHBfbWV0YWRhdGEgY3JlYXRlOnVzZXJzX2FwcF9tZXRhZGF0YSBjcmVhdGU6dXNlcl90aWNrZXRzIHJlYWQ6Y2xpZW50cyB1cGRhdGU6Y2xpZW50cyBkZWxldGU6Y2xpZW50cyBjcmVhdGU6Y2xpZW50cyByZWFkOmNsaWVudF9rZXlzIHVwZGF0ZTpjbGllbnRfa2V5cyBkZWxldGU6Y2xpZW50X2tleXMgY3JlYXRlOmNsaWVudF9rZXlzIHJlYWQ6Y29ubmVjdGlvbnMgdXBkYXRlOmNvbm5lY3Rpb25zIGRlbGV0ZTpjb25uZWN0aW9ucyBjcmVhdGU6Y29ubmVjdGlvbnMgcmVhZDpyZXNvdXJjZV9zZXJ2ZXJzIHVwZGF0ZTpyZXNvdXJjZV9zZXJ2ZXJzIGRlbGV0ZTpyZXNvdXJjZV9zZXJ2ZXJzIGNyZWF0ZTpyZXNvdXJjZV9zZXJ2ZXJzIHJlYWQ6ZGV2aWNlX2NyZWRlbnRpYWxzIHVwZGF0ZTpkZXZpY2VfY3JlZGVudGlhbHMgZGVsZXRlOmRldmljZV9jcmVkZW50aWFscyBjcmVhdGU6ZGV2aWNlX2NyZWRlbnRpYWxzIHJlYWQ6cnVsZXMgdXBkYXRlOnJ1bGVzIGRlbGV0ZTpydWxlcyBjcmVhdGU6cnVsZXMgcmVhZDplbWFpbF9wcm92aWRlciB1cGRhdGU6ZW1haWxfcHJvdmlkZXIgZGVsZXRlOmVtYWlsX3Byb3ZpZGVyIGNyZWF0ZTplbWFpbF9wcm92aWRlciBibGFja2xpc3Q6dG9rZW5zIHJlYWQ6c3RhdHMgcmVhZDp0ZW5hbnRfc2V0dGluZ3MgdXBkYXRlOnRlbmFudF9zZXR0aW5ncyByZWFkOmxvZ3MgcmVhZDpzaGllbGRzIGNyZWF0ZTpzaGllbGRzIGRlbGV0ZTpzaGllbGRzIHJlYWQ6Z3JhbnRzIGRlbGV0ZTpncmFudHMgcmVhZDpndWFyZGlhbl9mYWN0b3JzIHVwZGF0ZTpndWFyZGlhbl9mYWN0b3JzIHJlYWQ6Z3VhcmRpYW5fZW5yb2xsbWVudHMgZGVsZXRlOmd1YXJkaWFuX2Vucm9sbG1lbnRzIGNyZWF0ZTpndWFyZGlhbl9lbnJvbGxtZW50X3RpY2tldHMgcmVhZDp1c2VyX2lkcF90b2tlbnMifQ.oqheDADiW3ueJZhoP3LtODWqWYGsNEziRzHP6ASOhxzQATcREp4fqGXx2I2yNjYFSxRhPqKofYSeaMFsmgegvWvKmsonYbjhYDF8T0DIowSbE2beXmPb38puEZ3Ij4isLLQlp_1qAy7YGYvmrHJnnPcvhZGD7MB9o31Sw6vqK3jRG5KKfzT-PfqSsY2qjZRkmFtS2GjsmtGfs3UZy6RmDGH1RYnmwNRpYggrvTsLscVeW_KEUjbq68IB2Bv8Q3Wv7lzQ7AuPFFWXzNZJrx3MPDeMoGfALXjbo1lUonomoAMRV2fX1N-_JVYo2SFuhQBc9YN27vV7_DbhyPE99LL5Zg"

//MemView - aut0 api token
type MemView struct {
	Comments []database.Comment
	Mem      database.Mem
}

//Activity - aut0 api token
type Activity struct {
	MemID       int
	Description string
	DateTime    string
}

func init() {
	fmt.Println("Init package handlers")
}

func getNickname(userID string) string {
	url := "https://k4liber.eu.auth0.com/api/v2/users/" + userID
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", "Bearer "+APIToken)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var raw map[string]string
	json.Unmarshal(body, &raw)
	fmt.Println(raw)
	if raw["nickname"] != "" {
		return raw["nickname"]
	}
	return raw["name"]
}

func setHeaders(w http.ResponseWriter, req *http.Request) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	var origin = req.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, nickname")
	w.Header().Set("Allow", "*")
	return w
}

func uploadUserComments(nickname string, avatarName string) {
	var photoURL = HostName + "/resources/avatars/" + avatarName
	database.UserCommentUpdate(nickname, photoURL)
}

/*IconHandler ...
... */
func IconHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
	return http.HandlerFunc(fn)
}

/*IndexHandler ...
... */
func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
	return http.HandlerFunc(fn)
}

/*MemsHandler ...
... */
var MemsHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w = setHeaders(w, req)
	if req.Method == "OPTIONS" {
		return
	}
	payload, _ := json.Marshal(database.GetMems("janbielecki94"))
	w.Write([]byte(payload))
})

/*MemHandler ...
... */
var MemHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w = setHeaders(w, req)
	if req.Method == "OPTIONS" {
		return
	}
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Error in MemHandler!")
	}
	mem := database.GetMem(id, req.Header.Get("nickname"))
	comments := database.GetComments(id, req.Header.Get("nickname"))
	memView := MemView{
		Comments: comments,
		Mem:      mem,
	}
	payload, _ := json.Marshal(memView)

	w.Write([]byte(payload))
})

/*ProfileHandler ...
... */
var ProfileHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w = setHeaders(w, req)
	if req.Method == "OPTIONS" {
		return
	}
	vars := mux.Vars(req)
	nickname := vars["nickname"]
	payload, _ := json.Marshal(database.GetProfileMems(nickname))

	w.Write([]byte(payload))
})

/*ActivitiesHandler ...
... */
var ActivitiesHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w = setHeaders(w, req)
	if req.Method == "OPTIONS" {
		return
	}
	vars := mux.Vars(req)
	nickname := vars["nickname"]
	var mems = database.GetProfileMems(nickname)
	var comments = database.GetProfileComments(nickname)
	var commentLikes = database.GetProfileCommentLike(nickname)
	var memLikes = database.GetProfileMemLike(nickname)

	var activities []Activity
	var activity Activity
	var signature = ""
	for _, mem := range mems {
		if len(mem.Signature) > 30 {
			signature = mem.Signature[0:29] + "..."
		} else {
			signature = mem.Signature
		}
		activity.Description = "Your article '" + signature + "' has been added."
		activity.MemID = mem.ID
		activity.DateTime = mem.DateTime
		activities = append(activities, activity)
	}
	for _, comment := range comments {
		if len(comment.Content) > 30 {
			signature = comment.Content[0:29] + "..."
		} else {
			signature = comment.Content
		}
		activity.Description = "Your comment '" + signature + "' has been added."
		activity.MemID = comment.MemID
		activity.DateTime = comment.DateTime
		activities = append(activities, activity)
	}
	for _, commentLike := range commentLikes {
		activity.Description = "Your like this comment!"
		activity.MemID = commentLike.MemID
		activity.DateTime = commentLike.DateTime
		activities = append(activities, activity)
	}
	for _, memLike := range memLikes {
		activity.Description = "Your like this article!"
		activity.MemID = memLike.MemID
		activity.DateTime = memLike.DateTime
		activities = append(activities, activity)
	}
	payload, _ := json.Marshal(activities)
	w.Write([]byte(payload))
})

/*CategoryHandler ...
... */
var CategoryHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w = setHeaders(w, req)
	if req.Method == "OPTIONS" {
		return
	}
	vars := mux.Vars(req)
	category := vars["category"]
	nickname := req.Header.Get("nickname")
	payload, _ := json.Marshal(database.GetCategoryMems(category, nickname))
	w.Write([]byte(payload))
})

/*PreHandler ...
... */
var PreHandler = func(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var nickname = getNickname(req.FormValue("userID"))
		var authorNickname = req.FormValue("authorNickname")
		if nickname != authorNickname {
			fmt.Println(authorNickname)
			fmt.Println(nickname)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Cannot authorize!"))
			return
		}
		handler(w, req)
	}
}

/*AddMemHandler ...
... */
var AddMemHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	dateTime := time.Now().Format(time.RFC3339)
	var title = req.FormValue("title")
	var imgExt = "." + req.FormValue("extension")
	var category = req.FormValue("category")
	var authorNickname = req.FormValue("authorNickname")
	var content = req.FormValue("comment")
	var authorPhoto = req.FormValue("profilePicture")

	var mem database.MemDB
	mem.Signature = title
	mem.ImgExt = imgExt
	mem.DateTime = dateTime
	mem.AuthorNickname = authorNickname
	mem.Category = category
	var memID = database.InsertMem(mem)

	var comment database.CommentDB
	comment.MemID = memID
	comment.AuthorNickname = authorNickname
	comment.AuthorPhoto = authorPhoto
	comment.Content = content
	comment.DateTime = dateTime
	database.InsertComment(comment)

	//Zapisanie zdjecia
	req.ParseMultipartForm(32 << 20)
	file, handler, err := req.FormFile("file")
	var success = true
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	var fileName = "./resources/mems/" + strconv.Itoa(memID) + imgExt
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	//Wysylanie odpowiedzi
	payload, _ := json.Marshal(success)
	w.Write([]byte(payload))
})

/*UploadAvatarHandler ...
... */
var UploadAvatarHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var nickname = req.FormValue("authorNickname")
	var extension = req.FormValue("extension")
	var avatarName = nickname + "." + extension
	//Zapisanie zdjecia
	req.ParseMultipartForm(32 << 20)
	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	var fileName = "./resources/avatars/" + avatarName
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	uploadUserComments(nickname, avatarName)

	//Wysylanie odpowiedzi
	payload, _ := json.Marshal(avatarName)
	w.Header().Set("Avatar-Name", avatarName)
	w.Write([]byte(payload))
})

/*AddCommentHandler ...
... */
var AddCommentHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	dateTime := time.Now().Format(time.RFC3339)
	var memID, _ = strconv.Atoi(req.FormValue("memID"))
	var nickname = req.FormValue("authorNickname")
	var profilePicture = req.FormValue("profilePicture")
	var content = req.FormValue("comment")

	var comment database.CommentDB
	comment.MemID = memID
	comment.AuthorNickname = nickname
	comment.AuthorPhoto = profilePicture
	comment.Content = content
	comment.DateTime = dateTime
	comment = database.InsertComment(comment)

	payload, _ := json.Marshal(comment)
	w.Write([]byte(payload))
})

/*AddMemPointHandler .
... */
var AddMemPointHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var success = true
	dateTime := time.Now().Format(time.RFC3339)
	var memID, _ = strconv.Atoi(req.FormValue("memID"))
	var authorNickname = req.FormValue("authorNickname")
	var memPoint database.MemPoint
	memPoint.AuthorNickname = authorNickname
	memPoint.MemID = memID
	memPoint.DateTime = dateTime
	database.AddMemPoint(memPoint)
	database.UpdateMemPoints(memID)
	//Wysylanie odpowiedzi
	payload, _ := json.Marshal(success)
	w.Write([]byte(payload))
})

func commentProjection(commentDB database.CommentDB) database.Comment {
	var comment database.Comment
	comment.ID = commentDB.ID
	comment.AuthorNickname = commentDB.AuthorNickname
	comment.AuthorPhoto = commentDB.AuthorPhoto
	comment.Content = commentDB.Content
	comment.DateTime = commentDB.DateTime
	comment.MemID = commentDB.MemID
	comment.Points = commentDB.Points
	return comment
}

/*AddCommentPointHandler .
... */
var AddCommentPointHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	dateTime := time.Now().Format(time.RFC3339)
	var memID, _ = strconv.Atoi(req.FormValue("memID"))
	var authorNickname = req.FormValue("authorNickname")
	var commentID, _ = strconv.Atoi(req.FormValue("commentID"))
	var commentPoint database.CommentPoint
	commentPoint.AuthorNickname = authorNickname
	commentPoint.MemID = memID
	commentPoint.DateTime = dateTime
	commentPoint.CommentID = commentID
	database.AddCommentPoint(commentPoint)
	var commentDB = database.UpdateCommentPoints(commentID)
	var comment = commentProjection(commentDB)
	comment.Like = true

	payload, _ := json.Marshal(comment)
	w.Write([]byte(payload))
})

/*DeleteMemPointHandler .
... */
var DeleteMemPointHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var success = true
	dateTime := time.Now().Format(time.RFC3339)
	var memID, _ = strconv.Atoi(req.FormValue("memID"))
	var authorNickname = req.FormValue("authorNickname")
	var memPoint database.MemPoint
	memPoint.AuthorNickname = authorNickname
	memPoint.MemID = memID
	memPoint.DateTime = dateTime
	database.DeleteMemPoint(memPoint)
	database.UpdateMemPoints(memID)
	//Wysylanie odpowiedzi
	payload, _ := json.Marshal(success)
	w.Write([]byte(payload))
})

/*DeleteCommentPointHandler .
... */
var DeleteCommentPointHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	dateTime := time.Now().Format(time.RFC3339)
	var commentID, _ = strconv.Atoi(req.FormValue("commentID"))
	var authorNickname = req.FormValue("authorNickname")
	var commentPoint database.CommentPoint
	commentPoint.AuthorNickname = authorNickname
	commentPoint.CommentID = commentID
	commentPoint.DateTime = dateTime
	database.DeleteCommentPoint(commentPoint)
	var commentDB = database.UpdateCommentPoints(commentID)
	var comment = commentProjection(commentDB)
	comment.Like = false

	payload, _ := json.Marshal(comment)
	w.Write([]byte(payload))
})

/*DeleteMemHandler .
... */
var DeleteMemHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var success = true
	var memID, _ = strconv.Atoi(req.FormValue("memID"))

	database.DeleteMemPoints(memID)
	database.DeleteComments(memID)
	database.DeleteCommentPoints(memID)
	database.DeleteMem(memID)

	payload, _ := json.Marshal(success)
	w.Write([]byte(payload))
})

/*DeleteCommentHandler .
... */
var DeleteCommentHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var success = true
	var commentID, _ = strconv.Atoi(req.FormValue("commentID"))

	database.DeleteCommentPointsByCommID(commentID)
	database.DeleteComment(commentID)

	payload, _ := json.Marshal(success)
	w.Write([]byte(payload))
})

/*AdminDeleteCommentHandler .
... */
var AdminDeleteCommentHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var success = true
	var commentID, _ = strconv.Atoi(req.FormValue("commentID"))
	var authorNickname = req.FormValue("authorNickname")

	if authorNickname != "janbielecki94" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - You are not an admin!"))
		return
	}

	database.DeleteCommentPointsByCommID(commentID)
	database.DeleteComment(commentID)

	payload, _ := json.Marshal(success)
	w.Write([]byte(payload))
})

/*AdminDeleteMemHandler .
... */
var AdminDeleteMemHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	var success = true
	var memID, _ = strconv.Atoi(req.FormValue("memID"))
	var authorNickname = req.FormValue("authorNickname")

	if authorNickname != "janbielecki94" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - You are not an admin!"))
		return
	}

	database.DeleteMemPoints(memID)
	database.DeleteComments(memID)
	database.DeleteCommentPoints(memID)
	database.DeleteMem(memID)

	payload, _ := json.Marshal(success)
	w.Write([]byte(payload))
})
