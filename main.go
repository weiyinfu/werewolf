package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Room struct {
	lock    *sync.Mutex
	timer   *time.Timer
	Id      string            `json:"id"`
	Info    map[string]int    `json:"info"`
	Manager string            `json:"manager"`
	Turn    int               `json:"turn"`
	People  map[string]string `json:"people"`
}

func getRole(room *Room) string {
	//获取一个角色
	room.lock.Lock()
	defer func() {
		room.lock.Unlock()
	}()
	info := room.Info
	had := map[string]int{}
	for _, role := range room.People {
		had[role] += 1
	}
	var left []string
	for role, cnt := range info {
		if cnt < had[role] {
			panic("invalid state")
		}
		leftCount := cnt - had[role]
		for i := 0; i < leftCount; i++ {
			left = append(left, role)
		}
	}
	if len(left) == 0 {
		return ""
	}
	ind := rand.Intn(len(left))
	return left[ind]
}
func getUser(ctx *gin.Context) string {
	//从cookie中获取用户信息
	userId, err := ctx.Cookie("userid")
	if err != nil {
		userId = randRoomId()
		log.Println("为用户生成ID", userId, err)
		ctx.SetCookie("userid", userId, 1<<30, "", "", false, true)
	}
	return userId
}
func getRoom(roomId string) *Room {
	//从游戏表中获取游戏房间
	room_, ok := games.Load(roomId)
	if !ok {
		return nil
	}
	room := room_.(*Room)
	return room
}
func randRoomId() string {
	//随机一个房间id
	return fmt.Sprintf("%v", time.Now().UnixNano()^int64(rand.Int31()))
}

var games sync.Map

const ROOM_TIMEOUT = 3600 * 2 * time.Second //房间如果两个小时没有活跃则删除之

func resetRoomTimer(room *Room) {
	//检查房间是否过期，删除旧的定时任务，添加新的定时任务
	if !room.timer.Reset(ROOM_TIMEOUT) {
		log.Println("重置时钟错误")
	}
}

func getView(room *Room, userId string) gin.H {
	if room.Manager == userId {
		return gin.H{
			"fetched": len(room.People),
			"manager": room.Manager,
			"info":    room.Info,
			"turn":    room.Turn,
		}
	} else {
		return gin.H{
			"fetched": len(room.People),
			"role":    room.People[userId],
			"info":    room.Info,
			"turn":    room.Turn,
		}
	}
}

func main() {
	games = sync.Map{}
	var x = gin.Default()
	x.GET("/test", func(context *gin.Context) {
		_, _ = context.Writer.WriteString(fmt.Sprintf("狼人杀发牌助手 userId=%v", getUser(context)))
	})
	x.POST("/api/create_room", func(context *gin.Context) {
		room := Room{}
		room.Id = randRoomId()
		type createForm struct {
			Game map[string]int `json:"game"`
		}
		form := &createForm{}
		err := context.BindJSON(&form)
		if err != nil {
			log.Println("创建房间绑定参数错误", err)
			context.String(http.StatusBadRequest, "创建房间绑定参数错误")
			return
		}
		room.Info = form.Game
		room.Manager = getUser(context)
		room.Turn = 1
		room.People = map[string]string{}
		room.lock = &sync.Mutex{}
		room.timer = time.AfterFunc(ROOM_TIMEOUT, func() {
			log.Println("删除房间", room.Id)
			games.Delete(room.Id)
		})
		games.Store(room.Id, &room)
		resetRoomTimer(&room)
		log.Println("创建房间成功", room)
		context.JSON(http.StatusOK, gin.H{
			"id": room.Id,
		})
	})
	x.GET("/api/newgame", func(context *gin.Context) {
		userId := getUser(context)
		roomId := context.Query("room")
		log.Println("new game for room", roomId)
		room := getRoom(roomId)
		resetRoomTimer(room)
		if room.Manager != userId {
			context.JSON(http.StatusForbidden, map[string]string{
				"info": "只有管理员才能开始新游戏",
			})
			return
		}
		room.lock.Lock()
		defer func() {
			room.lock.Unlock()
		}()
		room.Turn += 1
		room.People = map[string]string{}
		context.JSON(http.StatusOK, getView(room, userId))
	})
	x.GET("/api/fetch", func(context *gin.Context) {
		roomId := context.Query("room")
		log.Println("fetch room", roomId)
		room := getRoom(roomId)
		if room == nil {
			log.Println("no such room", roomId)
			context.String(http.StatusOK, "no such room")
			return
		}
		resetRoomTimer(room)
		userId := getUser(context)
		if userId == room.Manager {
			context.JSON(http.StatusOK, getView(room, userId))
			return
		}
		if _, ok := room.People[userId]; !ok {
			room.People[userId] = getRole(room)
		}
		context.JSON(http.StatusOK, getView(room, userId))
	})
	x.StaticFS("/front", http.Dir("dist/"))
	ipport := "0.0.0.0:9968"
	if os.Getenv("GIN_MODE") != "release" {
		ipport = "localhost:9968"
	}
	_ = x.Run(ipport)
}
