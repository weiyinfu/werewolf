package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Room struct {
	Id      string            `json:"id"`
	Info    map[string]int    `json:"info"`
	Manager string            `json:"manager"`
	Turn    int               `json:"turn"`
	People  map[string]string `json:"people"`
	Fetched int               `json:"fetched"`
}
type RuntimeRoom struct {
	lock  *sync.Mutex
	timer *time.Timer
	room  *Room
}

func getRole(room *RuntimeRoom) string {
	//获取一个角色
	room.lock.Lock()
	defer func() {
		room.lock.Unlock()
	}()
	info := room.room.Info
	had := map[string]int{}
	for _, role := range room.room.People {
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
		log.Println("发现未登录用户")
		panic("发现未登录用户")
	}
	return userId
}
func getRoom(roomId string) *RuntimeRoom {
	//从游戏表中获取游戏房间
	room_, ok := games.Load(roomId)
	if !ok {
		return nil
	}
	room := room_.(*RuntimeRoom)
	return room
}
func randRoomId() string {
	//随机一个房间id
	return fmt.Sprintf("%v", time.Now().UnixNano()^int64(rand.Int31()))
}

var games sync.Map

const ROOM_TIMEOUT = 3600 * 2 * time.Second //房间如果两个小时没有活跃则删除之

func resetRoomTimer(room *RuntimeRoom) {
	//检查房间是否过期，删除旧的定时任务，添加新的定时任务
	if !room.timer.Reset(ROOM_TIMEOUT) {
		log.Println("重置时钟错误")
	}
}
func main() {
	games = sync.Map{}
	var x = gin.Default()
	x.GET("/test", func(context *gin.Context) {
		_, _ = context.Writer.WriteString("狼人杀发牌助手")
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
		runtimeRoom := &RuntimeRoom{
			room: &room,
			lock: &sync.Mutex{},
		}
		runtimeRoom.timer = time.AfterFunc(ROOM_TIMEOUT, func() {
			log.Println("删除房间", room.Id)
			games.Delete(room.Id)
		})
		games.Store(room.Id, runtimeRoom)
		resetRoomTimer(runtimeRoom)
		log.Println("创建房间成功", room)
		context.JSON(http.StatusOK, room)
	})
	x.GET("/api/newgame", func(context *gin.Context) {
		userId := getUser(context)
		roomId := context.Query("room")
		log.Println("new game for room", roomId)
		room := getRoom(roomId)
		resetRoomTimer(room)
		if room.room.Manager != userId {
			context.JSON(http.StatusForbidden, map[string]string{
				"info": "只有管理员才能开始新游戏",
			})
			return
		}
		room.lock.Lock()
		defer func() {
			room.lock.Unlock()
		}()
		room.room.Turn += 1
		room.room.People = map[string]string{}
		room.room.Fetched = 0
		context.JSON(http.StatusOK, room.room)
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
		if userId == room.room.Manager {
			r := *room.room
			fmt.Println(r.Info)
			fmt.Println(room.room.Info)
			r.People = map[string]string{}
			ind := 1
			for _, role := range room.room.People {
				r.People[fmt.Sprint(ind)] = role
			}
			context.JSON(http.StatusOK, r)
			return
		}
		if _, ok := room.room.People[userId]; !ok {
			room.room.People[userId] = getRole(room)
			room.room.Fetched += 1
		}
		r := *room.room
		r.People = map[string]string{
			userId: room.room.People[userId],
		}
		r.Manager = "" //不能让普通用户看到房主的信息
		context.JSON(http.StatusOK, r)
	})
	x.StaticFS("/front", http.Dir("dist/"))
	_ = x.Run("localhost:9968")
}
