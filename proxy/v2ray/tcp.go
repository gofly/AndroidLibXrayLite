package v2ray

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/2dust/AndroidLibXrayLite/core"

	vnet "github.com/xtls/xray-core/common/net"
	vsession "github.com/xtls/xray-core/common/session"
	vcore "github.com/xtls/xray-core/core"
)

type tcpHandler struct {
	ctx context.Context
	v   *vcore.Instance
}

func (h *tcpHandler) handleInput(conn net.Conn, input io.ReadCloser) {
	defer func() {
		conn.Close()
		input.Close()
	}()
	io.Copy(conn, input)
}

func (h *tcpHandler) handleOutput(conn net.Conn, output io.WriteCloser) {
	defer func() {
		conn.Close()
		output.Close()
	}()
	io.Copy(output, conn)
}

func NewTCPHandler(ctx context.Context, instance *vcore.Instance) core.TCPConnHandler {
	return &tcpHandler{
		ctx: ctx,
		v:   instance,
	}
}

func (h *tcpHandler) Handle(conn net.Conn, target *net.TCPAddr) error {
	dest := vnet.DestinationFromAddr(target)
	sid := vsession.NewID()
	ctx := vsession.ContextWithID(h.ctx, sid)
	c, err := vcore.Dial(ctx, h.v, dest)
	if err != nil {
		return errors.New(fmt.Sprintf("dial V proxy connection failed: %v", err))
	}
	go h.handleInput(conn, c)
	go h.handleOutput(conn, c)
	log.Printf("new proxy connection for target: %s:%s", target.Network(), target.String())
	return nil
}
