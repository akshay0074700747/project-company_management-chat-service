package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akshay0074700747/project-and-company-management-chat-service/helpers"
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/usecase"
	"github.com/akshay0074700747/project-and-company-management-chat-service/internal/usecase/chat"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/companypb"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/projectpb"
	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"github.com/gorilla/websocket"
)

type ChatHandlers struct {
	Usecase     usecase.ChatUseaseInterface
	Upgrader    websocket.Upgrader
	CompanyConn companypb.CompanyServiceClient
	ProjectConn projectpb.ProjectServiceClient
	UserConn    userpb.UserServiceClient
}

func NewChatHandlers(usecase usecase.ChatUseaseInterface, compaddr, projectAddr, userAddr string) *ChatHandlers {
	compRes, _ := helpers.DialGrpc(compaddr)
	prijectRes, _ := helpers.DialGrpc(projectAddr)
	userRes, _ := helpers.DialGrpc(userAddr)
	return &ChatHandlers{
		Usecase: usecase,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		CompanyConn: companypb.NewCompanyServiceClient(compRes),
		ProjectConn: projectpb.NewProjectServiceClient(prijectRes),
		UserConn:    userpb.NewUserServiceClient(userRes),
	}
}

func (chat *ChatHandlers) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", chat.handler)
	http.ListenAndServe(":50006", mux)
}

func (chatt *ChatHandlers) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=======================here===============================")
	r.Header.Del("Sec-WebSocket-Extensions")
	userID := r.Header.Get("userID")
	projectID := r.Header.Get("projectID")
	companyID := r.Header.Get("companyID")

	if userID == "" {
		http.Error(w, "the userID cannot be empty", http.StatusBadRequest)
		return
	}

	var roomID string
	if projectID != "" {
		_, err := chatt.ProjectConn.IsMemberAccepted(context.TODO(), &projectpb.IsMemberAcceptedReq{
			UserID:    userID,
			ProjectID: projectID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			helpers.PrintErr(err, "error ahppened at IsMemberAccepted")
			return
		}
		roomID = projectID
	} else if companyID != "" {
		exists, err := chatt.CompanyConn.IsEmployeeExists(context.TODO(), &companypb.IsEmployeeExistsReq{
			CompanyID:  companyID,
			EmployeeID: userID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			helpers.PrintErr(err, "error ahppened at IsEmployeeExists")
			return
		}
		if exists.Exists {
			roomID = companyID
		} else {
			http.Error(w, "you are not a part of the company", http.StatusBadRequest)
			helpers.PrintErr(err, "error ahppened at IsEmployeeExists")
			return
		}
	} else {
		http.Error(w, "the projectid or companyID should be specified", http.StatusBadRequest)
		return
	}

	details, err := chatt.UserConn.GetUserDetails(context.TODO(), &userpb.GetUserDetailsReq{
		UserID: userID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		helpers.PrintErr(err, "error ahppened at GetUserDetails")
		return
	}

	conn, err := chatt.Upgrader.Upgrade(w, r, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		helpers.PrintErr(err, "error happened at upgrading the request")
		return
	}

	pool := chatt.Usecase.SpinupPoolifnotalreadyExists(roomID)

	client := chat.NewClient(conn, userID, details.Name, pool)

	client.Serve()
}
