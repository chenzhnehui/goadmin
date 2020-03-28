package utils

// +----------------------------------------------------------------------
// | GOadmin [ I CAN DO IT JUST IT ]
// +----------------------------------------------------------------------
// | Copyright (c) 2020~2030 http://www.woaishare.cn All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( http://www.apache.org/licenses/LICENSE-2.0 )
// +----------------------------------------------------------------------
// | Author: chenzhenhui <971688607@qq.com>
// +----------------------------------------------------------------------
// | 分享交流QQ群请加  1062428023
// +----------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"unsafe"
)

type Message struct {
	Message  string `json:"message"`
	Type     string `json:"type"`
	Clientid string `json:"clientid"`
}

//客户端
type clientOne struct {
	Conn    *websocket.Conn
	Session interface{}
	Send    chan Message
}

var (
	clientsList      = make(map[string]clientOne)       //client 数据
	clientUids       = make(map[string]string)          //client_id 到 uid的绑定关系,1对多
	groupsClient     = make(map[string]map[string]bool) //group_id 到 client_id的绑定关系,1对多
	Msg              = make(chan Message, 10)           // 消息通道，收到的消息
	HandshakeTimeout = 10                               //保持心跳时间
	Upgrader         = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: time.Duration(HandshakeTimeout) * time.Second, //心跳时间，最后一次发送数据超过该时间就断开
		CheckOrigin: func(r *http.Request) bool { //忽略跨越
			return true
		},
	}
)

/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func init() {
	wsping()
}

//断开链接，关系客户端
/**
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
 */
func close(client_id string, client ...interface{}) error {
	delete(clientsList, client_id)
	delete(clientUids, client_id) //删除uid绑定的客户端
	if len(client) > 0 {
		switch client[0].(type) {
		case clientOne:
			client[0].(clientOne).Conn.Close()
		}
	}
	for group_id, clients := range groupsClient {
		if _, ok := clients[client_id]; ok {
			delete(groupsClient[group_id], client_id)
			Msg <- Message{Clientid: client_id, Type: "levelgroup", Message: group_id}
		}
		if len(clients) == 0 {
			delete(groupsClient, group_id)
			Msg <- Message{Type: "ungroup", Message: group_id}
		}
	}
	return nil
}

/**
保持心跳
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func wsping() error {
	go func() {
		for {
			time.Sleep(time.Duration(HandshakeTimeout) * time.Second)
			WsSendToAll(Message{Type: "ping", Message: Now()})
		}
	}()
	return nil
}

/**
初始化连接
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func ConnectInit(ws *websocket.Conn) string {
	client := clientOne{Send: make(chan Message)}
	client.Conn = ws
	client_id := fmt.Sprintf("%p", unsafe.Pointer(ws))
	clientsList[client_id] = client

	go readMsg(&client)
	go writeMsg(&client)
	Msg <- Message{Clientid: client_id, Type: "join", Message: "用户加入"}
	return client_id
}

/**
在client通道发送消息
*/
func writeMsg(client *clientOne) {
	for {
		select {
		case data := <-client.Send:
			err := client.Conn.WriteJSON(data)
			if err != nil {
				client.Conn.Close()
				close(fmt.Sprintf("%p", unsafe.Pointer(client.Conn)))
			}
		}
	}
}

/**
在client通道读取消息
*/
func readMsg(client *clientOne) {
	client_id := fmt.Sprintf("%p", unsafe.Pointer(client.Conn))
	defer func() {
		//取出该clien_id绑定的group和uid
		message := map[string]interface{}{"uid": WsGetUidByClientId(client_id), "group": WsGetGroupCountByClientId(client_id)}
		Msg <- Message{Clientid: client_id, Type: "level", Message: JsonEncode(message)}
		client.Conn.Close()
		close(client_id)
	}()
	for {
		// 读取消息。如果连接断开，则会返回错误	// 由于WebSocket一旦连接，便可以保持长时间通讯，则该接口函数可以一直运行下去，直到连接断开
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil { // 如果返回错误，就退出循环
			break
		}
		msg.Clientid = client_id
		Msg <- msg
	}
}

/**
发送消息到所有客户端
clientids 不发消息的客户端
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsSendToAll(msg Message, clientids ...interface{}) {
	var noClientids []string
	if len(clientids) > 0 {
		for _, c := range clientids {
			noClientids = append(noClientids, c.(string))
		}
	}
	go func() {
		for clientid, client := range clientsList {
			if len(noClientids) > 0 && InArray(clientid, noClientids) {
				continue
			}
			client.Send <- msg
		}
	}()
}

/**
发送消息到客户端
client_id 客户端地址 字符串
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsSendToClient(client_id string, msg Message) error {
	if _, ok := clientsList[client_id]; ok {
		clientsList[client_id].Send <- msg
		return nil
	}
	return errors.New("客户端不存在")
}

/**
关闭客户端
client_id 客户端地址 字符串
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsCloseClient(client_id string) error {
	client, ok := clientsList[client_id]
	if !ok {
		return errors.New("客户端不存在")
	}
	//取出该clien_id绑定的group和uid
	message := map[string]interface{}{"uid": WsGetUidByClientId(client_id), "group": WsGetGroupCountByClientId(client_id)}
	Msg <- Message{Clientid: client_id, Type: "level", Message: JsonEncode(message)}
	return close(client_id, client)
}

/**
判断客户端是否在线
bool  true 在线 false 不在线
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsIsOnline(client_id string) bool {
	_, ok := clientsList[client_id]
	return ok
}

/**
将client_id与uid绑定，以便通过WsSendToUid($uid)发送数据，通过WsIsUidOnline($uid)用户是否在线。
uid解释：这里uid泛指用户id或者设备id，用来唯一确定一个客户端用户或者设备。
1、uid与client_id是一对多的关系，系统允许一个uid下有多个client_id。
2、但是一个client_id只能绑定一个uid，如果绑定多次uid，则只有最后一次绑定有效。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsBindUid(client_id, uid string) error {
	if _, ok := clientsList[client_id]; !ok {
		return errors.New("客户端不存在")
	}
	clientUids[client_id] = uid
	return nil
}

/**
将client_id与uid解绑。
注意：当client_id下线（连接断开）时会自动与uid解绑
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsUnBindUid(client_id, uid string) error {
	for c, u := range clientUids {
		if client_id == c && u == uid {
			delete(clientUids, client_id)
			break
		}
	}
	return nil
}

/**
uid绑定的client_id是否在线
判断$uid是否在线，此方法需要配合BindUid($client_uid, $uid)使用。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsIsUidOnline(uid string) bool {
	if len(WsGetClientIdByUid(uid)) > 0 {
		return true
	}
	return false
}

/**
返回一个数组，数组元素为与uid绑定的所有在线的client_id。如果没有在线的client_id则返回一个空数组。
此方法可以判断一个uid是否在线。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsGetClientIdByUid(uid string) []string {
	var clientids []string
	for client_id, cuid := range clientUids {
		if cuid == uid {
			if _, ok := clientsList[client_id]; ok {
				clientids = append(clientids, client_id)
			} else {
				delete(clientUids, client_id)
			}
		}
	}
	return clientids
}

/**
返回client_id绑定的uid，如果client_id没有绑定uid，则返回空。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsGetUidByClientId(client_id string) string {
	return clientUids[client_id]
}

/**
向uid绑定的所有在线client_id发送数据。
注意：默认uid与client_id是一对多的关系，如果当前uid下绑定了多个client_id，则多个client_id对应的客户端都会收到消息，这类似于PC QQ和手机QQ同时在线接收消息。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsSendToUid(uid string, msg Message) error {
	var err error
	for _, client_id := range WsGetClientIdByUid(uid) {
		err = WsSendToClient(client_id, msg)
	}
	return err
}

/**
将client_id加入某个组，以便通过sendToGroup发送数据。
1、同一个client_id可以加入多个分组，以便接收不同组发来的数据。
2、当client_id下线（连接断开）后，该client_id会自动从该分组中删除，开发者无需调用Gateway::leaveGroup。
3、如果对应分组的所有client_id都下线，则对应分组会被自动删除
client_id 加入群组
group_id 群组id
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsJoinGroup(client_id, group_id string) error {
	if group_id == "" {
		return errors.New("分组不存在")
	}
	if _, ok := clientsList[client_id]; !ok {
		return errors.New("当前用户不在线")
	}
	groupInfo := make(map[string]bool)
	if _, ok := groupsClient[group_id]; ok {
		groupInfo = groupsClient[group_id]
	}
	groupInfo[client_id] = true
	groupsClient[group_id] = groupInfo
	Msg <- Message{Clientid: client_id, Type: "joingroup", Message: group_id}
	return nil
}

/**
将client_id从某个组中删除，不再接收该分组广播(sendToGroup)发送的数据。
当client_id下线（连接断开）时，client_id会自动从它所属的各个分组中删除
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsLeaveGroup(client_id, group_id string) error {
	delete(groupsClient[group_id], client_id)
	if len(groupsClient[group_id]) == 0 {
		delete(groupsClient, group_id)
	}
	Msg <- Message{Clientid: client_id, Type: "levelgroup", Message: group_id}
	return nil
}

/**
取消分组，或者说解散分组。 取消分组后所有属于这个分组的用户的连接将被移出分组，此分组将不再存在
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsUngroup(group_id string) error {
	if _, ok := groupsClient[group_id]; ok {
		delete(groupsClient, group_id)
	}
	Msg <- Message{Type: "ungroup", Message: group_id}
	return nil
}

/**
发送消息到所有分组
clientids 不发消息的客户端
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsSendToGroup(msg Message, group_id string, clientids ...interface{}) error {
	groupInfo, ok := groupsClient[group_id]
	if !ok {
		return errors.New("分组不存在")
	}
	if len(groupInfo) == 0 {
		return errors.New("分组成员不存在")
	}
	var noClientids []string
	if len(clientids) > 0 {
		for _, c := range clientids {
			noClientids = append(noClientids, c.(string))
		}
	}
	go func() {
		for clientid, _ := range groupInfo {
			if len(noClientids) > 0 && InArray(clientid, noClientids) {
				continue
			}
			if _, ok := clientsList[clientid]; ok {
				clientsList[clientid].Send <- msg
			}
		}
	}()
	return nil
}

/**
获取某分组当前在线成连接数（多少client_id在线）。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsGetClientIdCountByGroup(group_id string) int {
	groupInfo, ok := groupsClient[group_id]
	if !ok {
		return 0
	}
	return len(groupInfo)
}

/**
通过client_id 获取该client_id所在分组。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsGetGroupCountByClientId(client_id string) []string {
	var groups []string
	for group_id, clients := range groupsClient {
		for clientid, _ := range clients {
			if clientid == client_id {
				groups = append(groups, group_id)
			}
		}
	}
	return groups
}

/**
获取当前在线连接总数（多少client_id在线）。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsGetAllClientIdCount() int {
	return len(clientsList)
}

/**
获取当前在线client_id信息。
* @Author  chenzhenhui <971688607@qq.com>
* @Copyright  2020~2030 http://www.woaishare.cn All rights reserved.
*/
func WsGetClientSessions(client_id string) interface{} {
	return clientsList[client_id].Session
}
